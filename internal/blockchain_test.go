package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockChainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockChain(RandomBlock(0))
	assert.Nil(t, err)

	return bc
}

func TestBlock_Hash(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
	assert.True(t, bc.HasBlock(0))
}

func TestBlockchain_AddBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	b := RandomBlockWithSignature(t, 1)
	assert.Nil(t, bc.AddBlock(b))

	for i := 2; i < 1000; i++ {
		block := RandomBlockWithSignature(t, uint32(i))
		err := bc.AddBlock(block)
		assert.Nil(t, err)
	}

	assert.Equal(t, uint32(999), bc.Height())
	assert.Equal(t, 1000, bc.len())

	assert.NotNil(t, bc.AddBlock(RandomBlockWithSignature(t, 100)))
}
