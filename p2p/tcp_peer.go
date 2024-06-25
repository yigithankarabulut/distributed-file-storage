package p2p

import (
	"net"
	"sync"
)

// TCPPeer represents a peer in a TCP network.
type TCPPeer struct {
	// the underlying connection of the peer. Which in this case
	// is a TCP connection.
	net.Conn

	// if we dial and retrieve a connection, outbound == true
	// if we accept and retrieve a connection, outbound == false
	outbound bool

	Wg *sync.WaitGroup
}

// TCPPeerOption is a functional option for configuring a TCPPeer.
type TCPPeerOption func(*TCPPeer)

// WithTCPPeerConn is a functional option for setting the connection of the TCPPeer.
func WithTCPPeerConn(conn net.Conn) TCPPeerOption {
	return func(p *TCPPeer) {
		p.Conn = conn
	}
}

// WithTCPPeerOutbound is a functional option for setting the outbound flag of the TCPPeer.
func WithTCPPeerOutbound(outbound bool) TCPPeerOption {
	return func(p *TCPPeer) {
		p.outbound = outbound
	}
}

// NewTCPPeer creates a new TCPPeer with the given options.
func NewTCPPeer(opts ...TCPPeerOption) *TCPPeer {
	p := &TCPPeer{
		Wg: &sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Send sends data to the peer.
// Implement the Peer interface.
func (p *TCPPeer) Send(data []byte) error {
	_, err := p.Conn.Write(data)
	return err
}
