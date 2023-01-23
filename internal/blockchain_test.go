package internal

import (
	"testing"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestBlock_Hash(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
	assert.True(t, bc.HasBlock(0))
}

func TestBlockchain_AddBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	for i := 1; i < 1000; i++ {
		block := RandomBlockWithSignature(t, uint32(i), getPrevBlockHash(t, bc, uint32(i-1)))
		err := bc.AddBlock(block)
		assert.Nil(t, err)
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}

	assert.Equal(t, uint32(999), bc.Height())
	assert.Equal(t, 1000, bc.len())

	assert.NotNil(t, bc.AddBlock(RandomBlockWithSignature(t, 100, types.Hash{})))
}

func newBlockChainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockChain(RandomBlock(0, types.Hash{}))
	assert.Nil(t, err)
	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}
