package service

import (
	"errors"

	"github.com/StepanchukYI/simple-blockchain/internal"
	"github.com/StepanchukYI/simple-blockchain/internal/types"
	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
	log "github.com/sirupsen/logrus"
)

type ServiceInterface interface {
	CreateNewBlock() (*internal.Block, error)
	ProcessTransaction(tx *internal.Transaction) error
	ProcessBlock(block *internal.Block) error
	ProcessGetBlocks(from, to uint32) (*BlocksMessage, error)
	ProcessBlocks(blocks []*internal.Block) error
	GetHeight() uint32
	GetCurrentHeader() (*internal.Header, error)
}

type Service struct {
	chain      *internal.Blockchain
	mempool    *internal.TxPool
	PrivateKey *edwards.PrivateKey
}

func NewService(privKey *edwards.PrivateKey) (*Service, error) {
	gen, err := genesisBlock()
	if err != nil {
		return nil, err
	}

	chain, err := internal.NewBlockChain(gen)
	if err != nil {
		return nil, err
	}

	service := &Service{
		chain:      chain,
		mempool:    internal.NewTxPool(1000),
		PrivateKey: privKey,
	}

	return service, nil
}

func (s *Service) CreateNewBlock() (*internal.Block, error) {
	currentHeader, err := s.GetCurrentHeader()
	if err != nil {
		return nil, err
	}

	txx := s.mempool.Pending()

	block, err := internal.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return nil, err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return nil, err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return nil, err
	}

	return block, nil
}

func (s *Service) ProcessTransaction(tx *internal.Transaction) error {
	hash := tx.Hash(internal.TxHasher{})

	if s.mempool.Contains(hash) {
		return nil
	}

	if !tx.Verify() {
		return errors.New("invalid transaction")
	}

	log.Info(
		"msg", "adding new tx to mempool",
		"hash", hash,
		"mempoolPending", s.mempool.PendingCount(),
	)

	s.mempool.Add(tx)

	return nil
}

func (s *Service) ProcessBlock(block *internal.Block) error {
	if err := s.chain.AddBlock(block); err != nil {
		log.Error("error", err)
		return err
	}
	return nil
}

func (s *Service) ProcessGetBlocks(from, to uint32) (*BlocksMessage, error) {
	var blocks = []*internal.Block{}
	var ourHeight = s.chain.Height()

	if to == 0 {
		for i := int(from); i <= int(ourHeight); i++ {
			block, err := s.chain.GetBlock(uint32(i))
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, block)
		}
	}

	blocksMsg := &BlocksMessage{
		Blocks: blocks,
	}

	return blocksMsg, nil
}

func (s *Service) ProcessBlocks(blocks []*internal.Block) error {
	for _, block := range blocks {
		if err := s.ProcessBlock(block); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetHeight() uint32 {
	return s.chain.Height()
}

func (s *Service) GetCurrentHeader() (*internal.Header, error) {
	return s.chain.GetHeader(s.chain.Height())
}

func genesisBlock() (*internal.Block, error) {
	header := &internal.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 000000,
	}

	b, _ := internal.NewBlock(header, nil)

	coinbase := edwards.PublicKey{}
	tx := internal.NewTransaction(nil)
	tx.From = coinbase
	tx.To = coinbase
	tx.Value = 10_000_000
	b.Transactions = append(b.Transactions, tx)

	privKey, err := edwards.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	if err := b.Sign(privKey); err != nil {
		return nil, err
	}

	return b, nil
}
