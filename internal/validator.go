package internal

import (
	"errors"
	"fmt"
)

var (
	ErrBlockAlreadyExist = errors.New("block already exist")
)

type Validator interface {
	ValidateBlock(block *Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (bv BlockValidator) ValidateBlock(block *Block) error {
	if bv.bc.HasBlock(block.Height) {
		return ErrBlockAlreadyExist
	}

	if block.Height != bv.bc.Height()+1 {
		return fmt.Errorf("block to high with hash (%s)", block.Hash(BlockHasher{}))
	}

	prevHeader, err := bv.bc.GetHeader(block.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)
	if hash != block.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", block.PrevBlockHash)
	}

	_, err = block.Verify()
	if err != nil {
		return errors.New("block not valid")
	}
	return nil
}
