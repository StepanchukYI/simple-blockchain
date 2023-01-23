package internal

import (
	"bytes"
	"testing"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_Sign(t *testing.T) {
	data := []byte("testData")
	tx := &Transaction{
		Data: data,
	}

	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
	assert.True(t, tx.Verify())

	otherPrivKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	tx.From = otherPrivKey.PublicKey()

	assert.False(t, tx.Verify())
}

func TestTransaction_Decode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))
	tx.hash = types.Hash{}

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}

func randomTxWithSignature(t *testing.T) *Transaction {
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	tx := Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(privKey))

	return &tx
}
