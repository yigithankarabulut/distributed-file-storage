package p2p

import "net"

const (
	// IncomingMessage is a constant that represents an incoming message.
	IncomingMessage = 0x1
	// IncomingStream is a constant that represents an incoming stream.
	IncomingStream = 0x2
)

// RPC holds any arbitrary data that is being sent over
// each transport between two nodes in the network.
type RPC struct {
	From    net.Addr
	Payload []byte
	Stream  bool
}
