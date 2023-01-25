package types

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

const (
	HashLen = 32
)

type Hash [HashLen]uint8

func (h Hash) IsZero() bool {
	for i := 0; i < HashLen; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func (h Hash) Bytes() []byte {
	b := make([]byte, HashLen)
	for i, hb := range h {
		b[i] = hb
	}
	return b
}

func (h Hash) ToSlice() []byte {
	b := make([]byte, HashLen)
	for i, hb := range h {
		b[i] = hb
	}
	return b
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

func HashFromBytes(b []byte) (Hash, error) {
	if len(b) != HashLen {
		return Hash{}, errors.New("byte len must be equal 32")
	}

	var value [HashLen]uint8
	for i, bt := range b {
		value[i] = bt
	}

	return value, nil
}

func RandomHash() Hash {
	b := make([]byte, HashLen)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	h, err := HashFromBytes(b)
	if err != nil {
		panic(err)
	}

	return h
}
