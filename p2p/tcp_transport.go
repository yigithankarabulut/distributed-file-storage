package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

// ErrNilListener is an error that is returned when the listener is nil.
var ErrNilListener = errors.New("listener is nil")

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
		rpcCh: make(chan RPC, 1024),
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

// Close implements the Transport interface, which will close the transport and stop
// listening for incoming connections.
func (t *TCPTransport) Close() error {
	if t.listener != nil {
		return t.listener.Close()
	}
	return ErrNilListener
}

// Dial implements the Transport interface, which will dial a connection to the given address
// and then handle the connection.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

// Addr implements the Transport interface, which will return the address of the transport.
func (t *TCPTransport) Addr() string {
	return t.ListenAddr
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

	log.Printf("TCP transport listening on port [%s]\n", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			log.Printf("tcp accept error: %s\n", err.Error())
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		log.Printf("dropping peer connection: %s", err)
		_ = conn.Close()
	}()

	peer := NewTCPPeer(
		WithTCPPeerConn(conn),
		WithTCPPeerOutbound(outbound),
	)
	if err = t.ShakeHands(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read Loop
	for {
		rpc := RPC{}
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr()

		if rpc.Stream {
			peer.wg.Add(1)
			fmt.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
			peer.wg.Wait()
			fmt.Printf("[%s] stream closed, resuming read loop\n", conn.RemoteAddr())
			continue
		}

		t.rpcCh <- rpc
	}
}
