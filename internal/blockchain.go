package internal

type Blockchain struct {
	store     Storage
	headers   []Header
	validator Validator
}

func NewBlockChain(block *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []Header{},
		store:   NewMemorystore(),
	}
	bc.SetValidator(NewBlockValidator(bc))

	err := bc.addBlockWithoutValidation(block)
	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) len() int {
	return len(bc.headers)
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.headers = append(bc.headers, b.Header)
	return bc.store.Put(b)
}
