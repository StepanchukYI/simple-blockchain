package types

import (
	"encoding/hex"
	"errors"
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
