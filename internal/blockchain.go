package internal

import (
	"fmt"
	"sync"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
	log "github.com/sirupsen/logrus"
)

type Blockchain struct {
	mux     sync.RWMutex
	store   Storage
	headers []*Header
	blocks  []*Block

	transactionsStore map[types.Hash]*Transaction
	blockStore        map[types.Hash]*Block

	validator Validator
}

func NewBlockChain(block *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers:           []*Header{},
		store:             NewMemorystore(),
		transactionsStore: make(map[types.Hash]*Transaction),
		blockStore:        make(map[types.Hash]*Block),
	}
	bc.SetValidator(NewBlockValidator(bc))

	err := bc.addBlockWithoutValidation(block)
	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) Height() uint32 {
	bc.mux.RLock()
	defer bc.mux.RUnlock()
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) len() int {
	bc.mux.RLock()
	defer bc.mux.RUnlock()
	return len(bc.headers)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.mux.RLock()
	defer bc.mux.RUnlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) GetBlock(height uint32) (*Block, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.mux.RLock()
	defer bc.mux.RUnlock()

	return bc.blocks[height], nil
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	bc.headers = append(bc.headers, b.Header)
	bc.blocks = append(bc.blocks, b)

	bc.blockStore[b.Hash(BlockHasher{})] = b

	for _, tx := range b.Transactions {
		bc.transactionsStore[tx.Hash(TxHasher{})] = tx
	}

	log.Info(
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}

func (bc *Blockchain) GetTransactionByHash(hash types.Hash) (*Transaction, error) {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

	tx, ok := bc.transactionsStore[hash]
	if !ok {
		return nil, fmt.Errorf("could not find tx with hash (%s)", hash)
	}

	log.Info(
		"msg", "GetTransactionByHash",
		"hash", hash,
	)

	return tx, nil
}
