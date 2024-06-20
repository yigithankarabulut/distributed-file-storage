package p2p

import (
	"log"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

type TCPTransportOption func(*TCPTransport)

func WithListenAddr(addr string) TCPTransportOption {
	return func(t *TCPTransport) {
		t.listenAddr = addr
	}
}

func NewTCPTransport(opts ...TCPTransportOption) *TCPTransport {
	t := &TCPTransport{
		shakeHands: NOPHandshakeFunc,
		peers:      make(map[net.Addr]Peer),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *TCPTransport) ListenAndAccept() error {
	var (
		err error
	)
	t.listener, err = net.Listen("tcp", t.listenAddr)
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

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(
		WithTCPPeerConn(conn),
		WithTCPPeerOutbound(true),
	)
	if err := t.shakeHands(peer); err != nil {

	}

	// Read Loop
	msg := &Temp{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			log.Printf("tcp decode error: %s\n", err.Error())
			continue
		}
	}

}
