package main

import (
	"encoding/hex"
	"fmt"

	"github.com/StepanchukYI/simple-blockchain/pkg"
)

func main() {
	key, err := pkg.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(key.Bytes()))
}
