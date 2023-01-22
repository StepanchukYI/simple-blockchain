package internal

import (
	"testing"

	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_Sign(t *testing.T) {
	data := []byte("testData")
	tx := &Transaction{
		Data: data,
	}

	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
	assert.True(t, tx.Verify())

	otherPrivKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)
	tx.PublicKey = otherPrivKey.PublicKey()

	assert.False(t, tx.Verify())
}
