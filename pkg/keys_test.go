package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, len(privateKey.Bytes()), privateKeyLen)

	publicKey := privateKey.Public()
	assert.Equal(t, len(publicKey.Bytes()), publicKeyLen)
}

func TestGeneratePrivateKeyFromString(t *testing.T) {
	seed := "2200e39b93dcecc3e88eba779e2f118c3cbc48168eeef1a397d34afc6dc3bb1b"
	addressStr := "a3ff8ef6fe1131c2c2734cc281218d6555a8c0c4"
	privateKey, err := GeneratePrivateKeyFromString(seed)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(privateKey.Bytes()), privateKeyLen)
	address := privateKey.Public().Address()
	assert.Equal(t, address.String(), addressStr)
}

func TestPrivateKey_Sign(t *testing.T) {
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}
	publicKey := privateKey.Public()
	msg := []byte("message to sign")
	if err != nil {
		t.Error(err)
	}

	signed := privateKey.Sign(msg)
	assert.True(t, signed.Verify(publicKey, msg))

	assert.False(t, signed.Verify(publicKey, []byte("not same message to sign")))

	pvKey, err := GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}
	pbKey := pvKey.Public()
	assert.False(t, signed.Verify(pbKey, msg))
	assert.False(t, signed.Verify(pbKey, []byte("not same message to sign")))
}

func TestPublicKey_Address(t *testing.T) {
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}
	publicKey := privateKey.Public()
	address := publicKey.Address()
	assert.Equal(t, addressLen, len(address.Bytes()))

}
