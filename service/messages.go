package service

import "github.com/StepanchukYI/simple-blockchain/internal"

type BlocksMessage struct {
	Blocks []*internal.Block
}
