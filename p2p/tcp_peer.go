package p2p

import "net"

// TCPPeer represents a peer in a TCP network.
type TCPPeer struct {
	// conn is the underlying connection to the peer.
	conn net.Conn

	// if we dial and retrieve a connection, outbound == true
	// if we accept and retrieve a connection, outbound == false
	outbound bool
}

// TCPPeerOption is a functional option for configuring a TCPPeer.
type TCPPeerOption func(*TCPPeer)

// WithTCPPeerConn is a functional option for setting the connection of the TCPPeer.
func WithTCPPeerConn(conn net.Conn) TCPPeerOption {
	return func(p *TCPPeer) {
		p.conn = conn
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
	p := &TCPPeer{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Close closes the underlying connection to the peer.
// Implement the Peer interface.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}
