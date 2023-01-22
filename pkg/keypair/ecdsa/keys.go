package edwards

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"math/big"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
)

const (
	privateKeyLen = 64
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

type PublicKey []byte

type Signature struct {
	S *big.Int
	R *big.Int
}

func NewPrivateKeyFromReader(r io.Reader) PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), r)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

func GeneratePrivateKey() (PrivateKey, error) {
	return NewPrivateKeyFromReader(rand.Reader), nil
}

func (pv PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, pv.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func (pv PrivateKey) Public() PublicKey {
	return elliptic.MarshalCompressed(pv.key.PublicKey, pv.key.PublicKey.X, pv.key.PublicKey.Y)

}

func (pb PublicKey) Bytes() []byte {
	return pb
}

func (pb PublicKey) Address() types.Address {
	return types.AddressFromBytes(pb[len(pb)-types.AddressLen:])
}

func (s *Signature) Bytes() []byte {
	b := append(s.S.Bytes(), s.R.Bytes()...)
	return b
}

func (s *Signature) Verify(pubKey PublicKey, data []byte) bool {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pubKey)
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.Verify(key, data, s.R, s.S)
}
