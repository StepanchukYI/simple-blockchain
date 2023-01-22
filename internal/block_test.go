package internal

import (
	"testing"
	"time"

	"github.com/StepanchukYI/simple-blockchain/helpers"
	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
	"github.com/stretchr/testify/assert"
)

func RandomBlock(height uint32) *Block {
	header := Header{
		Version:       1,
		PrevBlockHash: helpers.RandomHash(),
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	tx := Transaction{
		Data: []byte("Test Data"),
	}

	return NewBlock(header, []Transaction{tx})
}

func RandomBlockWithSignature(t *testing.T, height uint32) *Block {
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	header := Header{
		Version:       1,
		PrevBlockHash: helpers.RandomHash(),
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	tx := Transaction{
		Data: []byte("Test Data"),
	}

	block := NewBlock(header, []Transaction{tx})

	err = block.Sign(privKey)
	assert.Nil(t, err)

	return block
}

func TestBlock_Sign(t *testing.T) {
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	b := RandomBlock(0)

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)

	assert.True(t, b.Verify())

	otherPrivKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	b.PublicKey = otherPrivKey.PublicKey()

	assert.False(t, b.Verify())

	b.Height = 100

	assert.False(t, b.Verify())
}
