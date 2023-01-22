package internal

import (
	"bytes"
	"encoding/gob"
	"io"

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

type Block struct {
	Header
	Transactions []Transaction
	PublicKey    edwards.PublicKey
	Signature    *edwards.Signature

	hash types.Hash
}

func NewBlock(h Header, tx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: tx,
	}
}

func (b *Block) Sign(privKey edwards.PrivateKey) error {
	sig, err := privKey.Sign(b.HeaderData())
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
	return b.Signature.Verify(b.PublicKey, b.HeaderData())
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Encode(w io.Writer, dec Encoder[*Block]) error {
	return dec.Encode(w, b)
}

func (b *Block) Hash(hashed Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hashed.Hash(b)
	}

	return b.hash
}

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(b.Header)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
