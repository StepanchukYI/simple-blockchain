package internal

import (
	"encoding/binary"
	"io"
)

type Transaction struct {
	Data []byte
}

func (tx *Transaction) EncodeBinary(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, &tx.Data); err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) DecodeBinary(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &tx.Data); err != nil {
		return err
	}

	return nil
}
