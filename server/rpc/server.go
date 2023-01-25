package rpc

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/StepanchukYI/simple-blockchain/internal"
	"github.com/StepanchukYI/simple-blockchain/service"
	log "github.com/sirupsen/logrus"
)

const defaultBlockTime = 5 * time.Second

type Config struct {
	Address string `mapstructure:"ADDRESS"  default:":8000"`
}

type ServerOpts struct {
	TCPTransport *TCPTransport
	Transports   []Transport
	ListenAddr   string
	ID           string

	RPCDecodeFunc RPCDecodeFunc
	BlockTime     time.Duration
	Service       service.ServiceInterface
}

type Server struct {
	TCPTransport *TCPTransport
	mux          sync.RWMutex
	ServerOpts

	peerCh  chan *TCPPeer
	peerMap map[net.Addr]*TCPPeer

	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
	txChan      chan *internal.Transaction
}

func NewServerOpts(cfg Config,
	id string,
	ts []Transport,
	service *service.Service,
) ServerOpts {
	opts := ServerOpts{
		ID:            id,
		BlockTime:     defaultBlockTime,
		RPCDecodeFunc: DefaultRPCDecodeFunc,
		ListenAddr:    cfg.Address,
		Transports:    ts,
		Service:       service,
	}

	return opts
}

func NewServer(opts ServerOpts) (*Server, error) {
	// Channel being used to communicate between the JSON RPC server
	// and the node that will process this message.
	txChan := make(chan *internal.Transaction)
	peerCh := make(chan *TCPPeer)

	tr := NewTCPTransport(opts.ListenAddr, peerCh)
	tr.peerCh = peerCh

	s := &Server{
		TCPTransport: tr,
		ServerOpts:   opts,
		peerCh:       peerCh,
		peerMap:      make(map[net.Addr]*TCPPeer),
		rpcCh:        make(chan RPC),
		quitCh:       make(chan struct{}, 1),
		txChan:       txChan,
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) Status() error {
	return nil
}

func (s *Server) Serve() error {

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				log.WithError(err)
				continue
			}

			if err := s.processMessage(msg); err != nil {
				if err != internal.ErrBlockAlreadyExist {
					log.WithError(err)
				}
			}
		case <-s.quitCh:
			break free

		}
	}

	return s.initTransport()
}

func (s *Server) initTransport() error {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				fmt.Printf("%+v", rpc)
			}
		}(tr)
	}

	return nil
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)

	log.Info("msg", "Starting validator loop", "blockTime", s.BlockTime)

	for {
		log.Info("creating new block")

		if err := s.createNewBlock(); err != nil {
			log.Error("create block error", err)
		}

		<-ticker.C
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.quitCh <- struct{}{}
	return nil
}

func (s *Server) requestBlocksLoop(peer net.Addr) error {
	ticker := time.NewTicker(3 * time.Second)

	for {
		ourHeight := s.Service.GetHeight()

		log.Info("msg", "requesting new blocks", "requesting height", ourHeight+1)

		// In this case we are 100% sure that the node has blocks heigher than us.
		getBlocksMessage := &GetBlocksMessage{
			From: ourHeight + 1,
			To:   0,
		}

		buf := new(bytes.Buffer)
		if err := gob.NewEncoder(buf).Encode(getBlocksMessage); err != nil {
			return err
		}

		s.mux.RLock()
		defer s.mux.RUnlock()

		msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())
		peer, ok := s.peerMap[peer]
		if !ok {
			return fmt.Errorf("peer %s not known", peer.conn.RemoteAddr())
		}

		if err := peer.Send(msg.Bytes()); err != nil {
			log.Error("error", "failed to send to peer", "err", err, "peer", peer)
		}

		<-ticker.C
	}
}

func (s *Server) broadcastTx(tx *internal.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(internal.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) broadcastBlock(b *internal.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(internal.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) broadcast(payload []byte) error {
	s.mux.RLock()
	defer s.mux.RUnlock()
	for netAddr, peer := range s.peerMap {
		if err := peer.Send(payload); err != nil {
			fmt.Printf("peer send error => addr %s [err: %s]\n", netAddr, err)
		}
	}

	return nil
}

func (s *Server) processTransaction(tx *internal.Transaction) error {
	if err := s.Service.ProcessTransaction(tx); err != nil {
		log.Error("error", err)
		return err
	}

	go func(tx *internal.Transaction) {
		err := s.broadcastTx(tx)
		if err != nil {
			log.Error("error", err)
		}
	}(tx)

	return nil
}

func (s *Server) processBlock(block *internal.Block) error {
	if err := s.Service.ProcessBlock(block); err != nil {
		log.Error("error", err)
		return err
	}

	go func(block *internal.Block) {
		err := s.broadcastBlock(block)
		if err != nil {
			log.Error("error", err)
		}
	}(block)

	return nil
}

func (s *Server) processGetBlocksMessage(from net.Addr, data *GetBlocksMessage) error {
	log.Info("msg", "received getBlocks message", "from", from)

	blocksMsg, err := s.Service.ProcessGetBlocks(data.From, data.To)
	if err != nil {
		log.Error("error", err)
		return err
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(blocksMsg); err != nil {
		return err
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	msg := NewMessage(MessageTypeBlocks, buf.Bytes())
	peer, ok := s.peerMap[from]
	if !ok {
		return fmt.Errorf("peer %s not known", peer.conn.RemoteAddr())
	}

	return peer.Send(msg.Bytes())
}

func (s *Server) processBlocksMessage(from net.Addr, data *BlocksMessage) error {
	log.Info("msg", "received BLOCKS!!!!!!!!", "from", from)

	err := s.Service.ProcessBlocks(data.Blocks)
	if err != nil {
		log.Info("error", err.Error())
		return err
	}

	return nil
}

func (s *Server) processStatusMessage(from net.Addr, data *StatusMessage) error {
	log.Info("msg", "received STATUS message", "from", from)

	if data.CurrentHeight <= s.Service.GetHeight() {
		log.Info("msg", "cannot sync blockHeight to low", "ourHeight", s.Service.GetHeight(), "theirHeight", data.CurrentHeight, "addr", from)
		return nil
	}

	go func(from net.Addr) {
		err := s.requestBlocksLoop(from)
		if err != nil {
			log.Error("error", err)
		}
	}(from)

	return nil
}

func (s *Server) processGetStatusMessage(from net.Addr, data *GetStatusMessage) error {
	log.Info("msg", "received getStatus message", "from", from)

	statusMessage := &StatusMessage{
		CurrentHeight: s.Service.GetHeight(),
		ID:            s.ID,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(statusMessage); err != nil {
		return err
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	peer, ok := s.peerMap[from]
	if !ok {
		return fmt.Errorf("peer %s not known", peer.conn.RemoteAddr())
	}

	msg := NewMessage(MessageTypeStatus, buf.Bytes())

	return peer.Send(msg.Bytes())
}

func (s *Server) processMessage(msg *DecodedMessage) (err error) {

	switch t := msg.Data.(type) {
	case *internal.Transaction:
		err = s.processTransaction(t)
		return err
	case *internal.Block:
		err = s.processBlock(t)
		return err
	case *GetBlocksMessage:
		return s.processGetBlocksMessage(msg.From, t)
	case *BlocksMessage:
		return s.processBlocksMessage(msg.From, t)
	case *GetStatusMessage:
		return s.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		return s.processStatusMessage(msg.From, t)
	}

	return err
}

func (s *Server) createNewBlock() error {
	block, err := s.Service.CreateNewBlock()
	if err != nil {
		return err
	}

	go func(block *internal.Block) {
		err := s.broadcastBlock(block)
		if err != nil {
			log.Error("error", err)
		}
	}(block)

	return nil
}
