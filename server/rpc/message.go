package rpc

import "github.com/StepanchukYI/simple-blockchain/internal"

type GetBlocksMessage struct {
	From uint32
	// If To is 0 the maximum blocks will be returned.
	To uint32
}

type BlocksMessage struct {
	Blocks []*internal.Block
}

type GetStatusMessage struct{}

type StatusMessage struct {
	// the id of the server
	ID            string
	Version       uint32
	CurrentHeight uint32
}
