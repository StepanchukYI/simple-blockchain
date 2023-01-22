package internal

import "errors"

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
		return errors.New("block already exist")
	}
	if !block.Verify() {
		return errors.New("block not valid")
	}
	return nil
}
