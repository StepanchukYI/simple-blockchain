package rpc

import (
	"net"
	"testing"

	"github.com/StepanchukYI/simple-blockchain/pkg/keypair/edwards"
	service "github.com/StepanchukYI/simple-blockchain/service"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t1Addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	assert.Nil(t, err)
	t1 := NewLocalTransport(t1Addr)
	privKey, err := edwards.GeneratePrivateKey()
	assert.Nil(t, err)

	serv, err := service.NewService(&privKey)
	assert.Nil(t, err)

	opts := NewServerOpts(Config{Address: "0.0.0.0"}, "TEST", []Transport{t1}, serv)
	s, err := NewServer(opts)
	assert.Nil(t, err)

	assert.NotNil(t, s)
}
