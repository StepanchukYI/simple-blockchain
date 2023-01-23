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

func TestAddBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	lenBlocks := 20
	for i := 1; i < lenBlocks; i++ {
		block := RandomBlock(t, uint32(i), getPrevBlockHash(t, bc, uint32(i)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks-1))
	assert.Equal(t, len(bc.headers), lenBlocks)
	assert.NotNil(t, bc.AddBlock(RandomBlock(t, 89, types.Hash{})))
}

func TestNewBlockchain(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(100))
}

func TestGetBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	lenBlocks := 20

	for i := 1; i < lenBlocks; i++ {
		block := RandomBlock(t, uint32(i), getPrevBlockHash(t, bc, uint32(i)))
		assert.Nil(t, bc.AddBlock(block))
		fetchedBlock, err := bc.GetBlock(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, fetchedBlock, block)
	}
}

func TestGetHeader(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	lenBlocks := 20

	for i := 1; i < lenBlocks; i++ {
		block := RandomBlock(t, uint32(i), getPrevBlockHash(t, bc, uint32(i)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(RandomBlock(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(RandomBlock(t, 3, types.Hash{})))
}

func newBlockChainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockChain(RandomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}
