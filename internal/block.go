package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
	Nonce         uint64
}

func (h Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}

type Block struct {
	*Header

	Transactions []*Transaction
	PublicKey    edwards.PublicKey
	Signature    *edwards.Signature

	hash types.Hash
}

func NewBlockFromPrevHeader(prevHeader *Header, tx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(tx)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		Height:        prevHeader.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     time.Now().UnixNano(),
	}

	return NewBlock(header, tx)
}

func NewBlock(h *Header, tx []*Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: tx,
	}, nil
}

func (b *Block) Sign(privKey edwards.PrivateKey) error {
	sig, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.PublicKey = privKey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() bool {
	if b.Signature == nil {
		return false
	}

	if b.Signature.Verify(b.PublicKey, b.Header.Bytes()) {

		return false
	}

	for _, tx := range b.Transactions {
		if !tx.Verify() {
			return false
		}
	}

	dataHash, _ := CalculateDataHash(b.Transactions)

	if dataHash != b.DataHash {
		return false
	}

	return b.Signature.Verify(b.PublicKey, b.Header.Bytes())
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(dec Encoder[*Block]) error {
	return dec.Encode(b)
}

func (b *Block) Hash(hashed Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hashed.Hash(b.Header)
	}

	return b.hash
}

func CalculateDataHash(txx []*Transaction) (hash types.Hash, err error) {
	buf := &bytes.Buffer{}

	for _, tx := range txx {
		if err = tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return
		}
	}

	hash = sha256.Sum256(buf.Bytes())

	return
}
