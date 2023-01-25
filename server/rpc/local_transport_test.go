package rpc

import (
	"bytes"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalTransport(t *testing.T) {
	t1Addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	assert.Nil(t, err)
	t1 := NewLocalTransport(t1Addr)

	t2Addr, err := net.ResolveTCPAddr("tcp", "localhost:8081")
	assert.Nil(t, err)
	t2 := NewLocalTransport(t2Addr)

	err = t1.Connect(t2)
	assert.Nil(t, err)

	payload := []byte("Hello, World!")

	// Test SendMessage
	err = t1.SendMessage(t2.Addr(), payload)
	assert.Nil(t, err)

	select {
	case rpc := <-t2.Consume():
		assert.Equal(t, rpc.From, t1.Addr())
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(rpc.Payload)
		assert.Nil(t, err)
		assert.Equal(t, buf.Bytes(), payload)
	default:
		t.Errorf("Expected message to be consumed")
	}

	// Test Broadcast
	t3Addr, _ := net.ResolveTCPAddr("tcp", "localhost:8082")
	t3 := NewLocalTransport(t3Addr)

	err = t1.Connect(t3)
	assert.Nil(t, err)
	err = t2.Connect(t3)
	assert.Nil(t, err)

	err = t1.Broadcast(payload)
	assert.Nil(t, err)

	select {
	case rpc := <-t3.Consume():
		assert.Equal(t, rpc.From, t1.Addr())
		assert.Nil(t, err)
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(rpc.Payload)
		assert.Nil(t, err)
		assert.Equal(t, buf.Bytes(), payload)
	default:
		t.Errorf("Expected message to be consumed")
	}
}
