package helpers

import (
	"crypto/rand"
	"io"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
)

func RandomHash() types.Hash {
	b := make([]byte, types.HashLen)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	h, err := types.HashFromBytes(b)
	if err != nil {
		panic(err)
	}

	return h
}
