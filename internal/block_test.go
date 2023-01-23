package internal

import (
	"testing"
	"time"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
	"github.com/stretchr/testify/assert"
)

func RandomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	tx := &Transaction{
		Data: []byte("Test Data"),
	}

	b, _ := NewBlock(header, []*Transaction{tx})
	return b
}

func RandomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	tx := &Transaction{
		Data: []byte("Test Data"),
	}

	block, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)

	err = block.Sign(privKey)
	assert.Nil(t, err)

	return block
}

func TestBlock_Sign(t *testing.T) {
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	b := RandomBlock(0, types.Hash{})

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
