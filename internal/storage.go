package internal

type Storage interface {
	Put(block *Block) error
	Get(height uint32) (Block, error)
}

type MemoryStore struct {
}

func NewMemorystore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Put(b *Block) error {
	return nil
}

func (s *MemoryStore) Get(height uint32) (Block, error) {
	return Block{}, nil
}
