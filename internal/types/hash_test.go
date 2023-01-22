package types

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomHash(t *testing.T) {
	hash, err := RandomHash()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(hash), hashLen)
}

func TestHashFromBytes(t *testing.T) {
	b := make([]byte, 20)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		t.Error(err)
	}

	h, err := HashFromBytes(b)
	assert.True(t, h.IsZero())
	assert.Error(t, err, "byte len must be equal 32")
}
