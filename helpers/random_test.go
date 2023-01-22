package helpers

import (
	"testing"

	"github.com/StepanchukYI/simple-blockchain/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRandomHash(t *testing.T) {
	hash := RandomHash()

	assert.Equal(t, len(hash), types.HashLen)
}
