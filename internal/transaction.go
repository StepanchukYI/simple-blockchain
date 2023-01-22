package internal

import (
	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
)

type Transaction struct {
	Data []byte

	PublicKey edwards.PublicKey
	Signature *edwards.Signature
}

func (tx *Transaction) Sign(privKey edwards.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.PublicKey = privKey.PublicKey()
	tx.Signature = sig

	return nil
}

func (tx *Transaction) Verify() bool {
	if tx.Signature == nil {
		return false
	}
	return tx.Signature.Verify(tx.PublicKey, tx.Data)
}
