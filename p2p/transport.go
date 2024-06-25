package p2p

import "net"

// Peer represents a connection to another node in the network.
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport represents a network transport that can listen for incoming
// connections and accept them, as well as send and receive messages.
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
