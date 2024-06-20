package main

import (
	"github.com/yigithankarabulut/distributed-file-storage/p2p"
	"log"
)

func main() {
	tr := p2p.NewTCPTransport(
		p2p.WithListenAddr(":4242"),
	)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("failed to listen and accept: %s\n", err.Error())
	}
	select {}
}
