package rpc

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t1Addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	assert.Nil(t, err)
	t1 := NewLocalTransport(t1Addr)

	opts := NewServerOpts(Config{Address: "0.0.0.0"}, "TEST", []Transport{t1})
	s, err := NewServer(opts)
	assert.Nil(t, err)

	assert.NotNil(t, s)
}
