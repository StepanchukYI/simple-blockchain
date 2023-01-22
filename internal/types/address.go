package types

import (
	"encoding/hex"
	"fmt"
)

const AddressLen = 20

type Address [AddressLen]uint8

func (a Address) Bytes() []byte {
	b := make([]byte, AddressLen)
	for i := 0; i < AddressLen; i++ {
		b[i] = a[i]
	}
	return b
}

func (a Address) ToSlice() []byte {
	b := make([]byte, AddressLen)
	for i := 0; i < AddressLen; i++ {
		b[i] = a[i]
	}
	return b
}

func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

func AddressFromBytes(b []byte) Address {
	if len(b) != AddressLen {
		msg := fmt.Sprintf("given bytes with length %d should be 20", len(b))
		panic(msg)
	}

	var value [AddressLen]uint8
	for i := 0; i < AddressLen; i++ {
		value[i] = b[i]
	}

	return Address(value)
}
