package edwards

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
)

const (
	privateKeyLen = 64
	publicKeyLen  = 32
	seedLen       = 32
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

type PublicKey ed25519.PublicKey

type Signature struct {
	Value []byte
}

func GeneratePrivateKeyFromString(s string) (PrivateKey, error) {
	seed, err := hex.DecodeString(s)
	if err != nil {
		return PrivateKey{}, err
	}

	return GeneratePrivateKeyFromSeed(seed)
}

func GeneratePrivateKeyFromSeed(seed []byte) (PrivateKey, error) {
	if len(seed) != seedLen {
		return PrivateKey{}, errors.New("invalid seed len")
	}

	return PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}, nil
}

func GeneratePrivateKey() (PrivateKey, error) {
	seed := make([]byte, seedLen)

	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}, nil
}

func (pv PrivateKey) Bytes() []byte {
	return pv.key
}

func (pv PrivateKey) Sign(data []byte) (*Signature, error) {
	return &Signature{
		Value: ed25519.Sign(pv.key, data),
	}, nil
}

func (pv PrivateKey) PublicKey() PublicKey {
	b := make([]byte, publicKeyLen)
	copy(b, pv.key[publicKeyLen:])

	return b
}

func (pb PublicKey) Address() types.Address {
	return types.AddressFromBytes(pb[publicKeyLen-types.AddressLen:])
}

func (s Signature) Bytes() []byte {
	return s.Value
}

func (s Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(pubKey), data, s.Value)
}
