package edwards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	signed, err := privateKey.Sign(msg)
	if err != nil {
		t.Error(err)
	}
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
