package internal

import (
	"bytes"
	"testing"
	"time"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
	"github.com/stretchr/testify/assert"
)

func RandomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	tx := randomTxWithSignature(t)
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	assert.Nil(t, b.Sign(privKey))

	return b
}

func RandomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey, err := edwards.GeneratePrivateKey()
	block := RandomBlock(t, height, prevBlockHash)
	assert.Nil(t, err)
	tx := randomTxWithSignature(t)
	block.AddTransaction(tx)
	err = block.Sign(privKey)
	assert.Nil(t, err)

	return block
}

func TestBlock_Sign(t *testing.T) {
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	b := RandomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	b := RandomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	ver, err := b.Verify()
	assert.True(t, ver)
	assert.Nil(t, err)

	otherPrivKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	b.PublicKey = otherPrivKey.PublicKey()
	ver, err = b.Verify()
	assert.False(t, ver)
	assert.NotNil(t, err)

	b.Height = 100
	ver, err = b.Verify()
	assert.False(t, ver)
	assert.NotNil(t, err)
}

func TestDecodeEncodeBlock(t *testing.T) {
	b := RandomBlock(t, 1, types.Hash{})
	buf := &bytes.Buffer{}
	assert.Nil(t, b.Encode(NewGobBlockEncoder(buf)))

	bDecode := new(Block)
	assert.Nil(t, bDecode.Decode(NewGobBlockDecoder(buf)))

	assert.Equal(t, b.Header, bDecode.Header)

	for i := 0; i < len(b.Transactions); i++ {
		b.Transactions[i].hash = types.Hash{}
		assert.Equal(t, b.Transactions[i], bDecode.Transactions[i])
	}

	assert.Equal(t, b.PublicKey, bDecode.PublicKey)
	assert.Equal(t, b.Signature, bDecode.Signature)
}
