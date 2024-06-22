package p2p

import (
	"errors"
	"log"
	"net"
)

// TCPTransport is a transport implementation that uses TCP as the underlying network protocol.
type TCPTransport struct {
	ListenAddr string
	ShakeHands HandshakeFunc
	Decoder    Decoder
	OnPeer     func(Peer) error

	listener net.Listener
	rpcCh    chan RPC
}

// TCPTransportOption is a functional option type for configuring a TCPTransport.
type TCPTransportOption func(*TCPTransport)

// WithListenAddr is a functional option for setting the listen address of the TCPTransport.
func WithListenAddr(addr string) TCPTransportOption {
	return func(t *TCPTransport) {
		t.ListenAddr = addr
	}
}

// WithHandshakeFunc is a functional option for setting the handshake function of the TCPTransport.
func WithHandshakeFunc(h HandshakeFunc) TCPTransportOption {
	return func(t *TCPTransport) {
		t.ShakeHands = h
	}
}

// WithDecoder is a functional option for setting the decoder of the TCPTransport.
func WithDecoder(d Decoder) TCPTransportOption {
	return func(t *TCPTransport) {
		t.Decoder = d
	}
}

// WithOnPeer is a functional option for setting the on peer function of the TCPTransport.
func WithOnPeer(f func(Peer) error) TCPTransportOption {
	return func(t *TCPTransport) {
		t.OnPeer = f
	}
}

// NewTCPTransport creates a new TCPTransport with the given options.
func NewTCPTransport(opts ...TCPTransportOption) *TCPTransport {
	t := &TCPTransport{
		rpcCh: make(chan RPC),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Consume implements the Transport interface, which will return a read-only channel
// for reading incoming messages received from another peer in the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

// ListenAndAccept implements the Transport interface, which will listen for incoming
// connections and accept them, and then handle the connection.
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		log.Printf("tcp listen error: %s\n", err.Error())
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			log.Printf("tcp accept error: %s\n", err.Error())
		}

		log.Printf("new incoming connection: %+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		log.Printf("dropping peer connection: %s", err)
		_ = conn.Close()
	}()

	peer := NewTCPPeer(
		WithTCPPeerConn(conn),
		WithTCPPeerOutbound(true),
	)
	if err = t.ShakeHands(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			log.Printf("tcp read error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcCh <- rpc
	}
}
