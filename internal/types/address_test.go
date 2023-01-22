package types

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressFromBytes(t *testing.T) {
	b := make([]byte, AddressLen)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		t.Error(err)
	}

	addr := AddressFromBytes(b)
	assert.Equal(t, len(addr.Bytes()), AddressLen)
	assert.Equal(t, len(addr.ToSlice()), AddressLen)
}
