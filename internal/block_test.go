package internal

import (
	"bytes"
	"testing"
	"time"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestHeader_EncodeBinary_DecodeBinary(t *testing.T) {
	hash, err := types.RandomHash()
	if err != nil {
		t.Error(err)
	}

	h := &Header{
		Version:   1,
		PrevBloc:  hash,
		Timestamp: time.Now().UnixNano(),
		Height:    2,
		Nonce:     454324,
	}

	buf := &bytes.Buffer{}
	err = h.EncodeBinary(buf)
	if err != nil {
		t.Error(err)
	}

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}

func TestBlock_EncodeBinary_DecodeBinary(t *testing.T) {
	hash, err := types.RandomHash()
	if err != nil {
		t.Error(err)
	}

	h := Header{
		Version:   1,
		PrevBloc:  hash,
		Timestamp: time.Now().UnixNano(),
		Height:    2,
		Nonce:     454324,
	}

	b := &Block{
		Header:       h,
		Transactions: nil,
	}

	buf := &bytes.Buffer{}
	err = b.EncodeBinary(buf)
	if err != nil {
		t.Error(err)
	}

	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, b, bDecode)
}

func TestBlock_Hash(t *testing.T) {
	hash, err := types.RandomHash()
	if err != nil {
		t.Error(err)
	}

	header := Header{
		Version:   1,
		PrevBloc:  hash,
		Timestamp: time.Now().UnixNano(),
		Height:    2,
		Nonce:     768945,
	}

	b := &Block{
		Header:       header,
		Transactions: nil,
	}

	h := b.Hash()
	assert.False(t, h.IsZero())
}
