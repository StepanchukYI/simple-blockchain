package pkg

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

const (
	privateKeyLen = 64
	publicKeyLen  = 32
	seedLen       = 32
	addressLen    = 20
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

type PublicKey struct {
	key ed25519.PublicKey
}

type Signature struct {
	value []byte
}

type Address struct {
	value []byte
}

func GeneratePrivateKeyFromString(s string) (*PrivateKey, error) {
	seed, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return GeneratePrivateKeyFromSeed(seed)
}

func GeneratePrivateKeyFromSeed(seed []byte) (*PrivateKey, error) {
	if len(seed) != seedLen {
		return nil, errors.New("invalid seed len")
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}, nil
}

func GeneratePrivateKey() (*PrivateKey, error) {
	seed := make([]byte, seedLen)

	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}, nil
}

func (pv *PrivateKey) Bytes() []byte {
	return pv.key
}

func (pv *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(pv.key, msg),
	}
}

func (pv *PrivateKey) Public() *PublicKey {
	b := make([]byte, publicKeyLen)
	copy(b, pv.key[publicKeyLen:])

	return &PublicKey{
		key: b,
	}
}

func (pb *PublicKey) Bytes() []byte {
	return pb.key
}

func (pb *PublicKey) Address() Address {
	return Address{
		value: pb.key[publicKeyLen-addressLen:],
	}
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(key *PublicKey, msg []byte) bool {
	return ed25519.Verify(key.key, msg, s.value)
}

func (a Address) Bytes() []byte {
	return a.value
}

func (a Address) String() string {
	return hex.EncodeToString(a.value)
}
