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

type TCPPeerOption func(*TCPPeer)

func WithTCPPeerConn(conn net.Conn) TCPPeerOption {
	return func(p *TCPPeer) {
		p.conn = conn
	}
}

func WithTCPPeerOutbound(outbound bool) TCPPeerOption {
	return func(p *TCPPeer) {
		p.outbound = outbound
	}
}

func NewTCPPeer(opts ...TCPPeerOption) *TCPPeer {
	p := &TCPPeer{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
