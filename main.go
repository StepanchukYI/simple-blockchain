package main

import (
	"encoding/hex"
	"fmt"

	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
)

func main() {
	key, err := edwards.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(key.Bytes()))
}
