package types

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashFromBytesError(t *testing.T) {
	b := make([]byte, 20)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		t.Error(err)
	}

	h, err := HashFromBytes(b)
	assert.Error(t, err, "byte len must be equal 32")
	assert.True(t, h.IsZero())
}

func TestHashFromBytes(t *testing.T) {
	b := make([]byte, HashLen)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		t.Error(err)
	}

	h, err := HashFromBytes(b)
	assert.Equal(t, len(h.Bytes()), HashLen)
	assert.Equal(t, len(h.ToSlice()), HashLen)
}

func TestRandomHash(t *testing.T) {
	hash := RandomHash()

	assert.Equal(t, len(hash), HashLen)
}
