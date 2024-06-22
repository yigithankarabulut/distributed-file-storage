package main

import (
	"log"

	"github.com/yigithankarabulut/distributed-file-storage/p2p"
)

func onPeer(peer p2p.Peer) error {
	// log.Println("doing some logic with the peer outside the TCPTransport")
	return peer.Close()
}

func main() {
	tr := p2p.NewTCPTransport(
		p2p.WithListenAddr(":4242"),
		p2p.WithHandshakeFunc(p2p.NOPHandshakeFunc),
		p2p.WithDecoder(&p2p.DefaultDecoder{}),
		p2p.WithOnPeer(onPeer),
	)

	go func() {
		for {
			msg := <-tr.Consume()
			log.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("failed to listen and accept: %s\n", err.Error())
	}
	select {}
}
