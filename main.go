package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/yigithankarabulut/distributed-file-storage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransport := p2p.NewTCPTransport(
		p2p.WithListenAddr(listenAddr),
		p2p.WithHandshakeFunc(p2p.NOPHandshakeFunc),
		p2p.WithDecoder(&p2p.DefaultDecoder{}),
	)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":4242", "")
	s2 := makeServer(":2424", ":4242")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(2 * time.Second)
	go func() {
		log.Fatal(s2.Start())
	}()
	time.Sleep(3 * time.Second)

	for i := 0; i < 10; i++ {
		data := bytes.NewReader([]byte("my big data file here!"))
		s2.Store(fmt.Sprintf("myprivatedata_%d", i), data)
		time.Sleep(5 * time.Millisecond)
	}

	// r, err := s2.Get("myprivatedata")
	// if err != nil {
	//	log.Fatal(err)
	//}
	//
	//b, err := io.ReadAll(r)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(string(b))

	select {}
}
