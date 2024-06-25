package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/yigithankarabulut/distributed-file-storage/fileserver"
	"github.com/yigithankarabulut/distributed-file-storage/p2p"
	"github.com/yigithankarabulut/distributed-file-storage/store"
)

func makeServer(listenAddr string, nodes ...string) *fileserver.FileServer {
	tcpTransport := p2p.NewTCPTransport(
		p2p.WithListenAddr(listenAddr),
		p2p.WithHandshakeFunc(p2p.NOPHandshakeFunc),
		p2p.WithDecoder(&p2p.DefaultDecoder{}),
	)

	fileServerOpts := fileserver.ServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: store.CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := fileserver.NewFileServer(fileServerOpts)

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

	// data := bytes.NewReader([]byte("my big data file here!"))
	// _ = s2.Store("coolPicture.jpg", data)
	// time.Sleep(5 * time.Millisecond)

	r, err := s2.Get("coolPicture.jpg")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

// 8:04. creating crypto file.
