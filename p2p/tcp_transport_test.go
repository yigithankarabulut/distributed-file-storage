package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4242"
	tr := NewTCPTransport(
		WithListenAddr(listenAddr),
		WithHandshakeFunc(NOPHandshakeFunc),
		WithDecoder(&DefaultDecoder{}),
	)

	assert.Equal(t, tr.ListenAddr, listenAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
