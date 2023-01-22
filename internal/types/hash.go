package types

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

const (
	hashLen = 32
)

type Hash [hashLen]uint8

func (h Hash) IsZero() bool {
	for i := 0; i < hashLen; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func (h Hash) ToSlice() []byte {
	b := make([]byte, hashLen)
	for i, hb := range h {
		b[i] = hb
	}
	return b
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

func HashFromBytes(b []byte) (Hash, error) {
	if len(b) != hashLen {
		return Hash{}, errors.New("byte len must be equal 32")
	}

	var value [hashLen]uint8
	for i, bt := range b {
		value[i] = bt
	}

	return value, nil
}

func RandomHash() (Hash, error) {
	b := make([]byte, hashLen)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return Hash{}, err
	}

	return HashFromBytes(b)
}
