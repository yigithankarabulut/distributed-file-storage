package fileserver

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/yigithankarabulut/distributed-file-storage/crypto"
	"github.com/yigithankarabulut/distributed-file-storage/p2p"
	"github.com/yigithankarabulut/distributed-file-storage/store"
)

// ServerOpts is a struct that contains the configuration for the file server.
type ServerOpts struct {
	ID                string
	EncryptKey        []byte
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc store.PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

// FileServer is a struct that contains the configuration for the file server.
type FileServer struct {
	ServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	Storage  *store.Store
	doneChan chan struct{}
}

// NewFileServer creates a new file server instance with the given options.
func NewFileServer(opts ServerOpts) *FileServer {
	s := store.NewStore(
		store.WithRoot(opts.StorageRoot),
		store.WithPathTransformFunc(opts.PathTransformFunc),
	)

	if len(opts.ID) == 0 {
		opts.ID = crypto.GenerateID()
	}

	return &FileServer{
		ServerOpts: opts,
		Storage:    s,
		doneChan:   make(chan struct{}),
		peers:      make(map[string]p2p.Peer),
	}
}

// Start starts the file server.
func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootstrapNodes) > 0 {
		s.bootstrapNetwork()
	}

	s.loop()

	return nil
}

// Stop stops the file server.
func (s *FileServer) Stop() {
	close(s.doneChan)
}

// OnPeer is a callback function that is called when a peer is connected to the file server.
func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote: %s\n", p.RemoteAddr())

	return nil
}

// Get gets the data from the file server.
// It reads the data from the store if it exists, otherwise it fetches the data from the network.
func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.Storage.Has(s.ID, key) {
		log.Printf("[%s] serving file (%s) from local disk\n", s.Transport.Addr(), key)
		_, r, err := s.Storage.Read(s.ID, key)
		return r, err
	}

	log.Printf("[%s] dont have the file (%s) locally, fetching from network...\n", s.Transport.Addr(), key)

	msg := Message{
		Payload: MessageGetFile{
			ID:  s.ID,
			Key: crypto.HashKey(key),
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 500)

	for _, peer := range s.peers {
		// First read the file size, so we can limit the amount of bytes that we read
		// from the connection, so it will not keep hanging.
		var fileSize int64
		if err := binary.Read(peer, binary.LittleEndian, &fileSize); err != nil {
			return nil, err
		}

		n, err := s.Storage.WriteDecrypt(
			s.EncryptKey,
			s.ID,
			key,
			io.LimitReader(peer, fileSize),
		)
		if err != nil {
			return nil, err
		}

		log.Printf("[%s] received (%d) bytes over the network from (%s)\n", s.Transport.Addr(), n, peer.RemoteAddr())

		peer.CloseStream()
	}

	_, r, err := s.Storage.Read(s.ID, key)
	return r, err
}

// Store stores the data in the file server.
// It writes the data to the store and then broadcasts the message to the peers.
func (s *FileServer) Store(key string, r io.Reader) error {
	var (
		fileBuffer = new(bytes.Buffer)
		tee        = io.TeeReader(r, fileBuffer)
	)

	size, err := s.Storage.Write(s.ID, key, tee)
	if err != nil {
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			ID:   s.ID,
			Key:  crypto.HashKey(key),
			Size: size + 16,
		},
	}

	if err = s.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 5)

	peers := make([]io.Writer, 0, len(s.peers))
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	if _, err = mw.Write([]byte{p2p.IncomingStream}); err != nil {
		return err
	}

	n, err := crypto.CopyEncrypt(s.EncryptKey, fileBuffer, mw)
	if err != nil {
		return err
	}

	log.Printf("[%s] received and written (%d) bytes to disk: ", s.Transport.Addr(), n)

	return nil
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		_ = peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Printf("file server stopped due to error or user quit action\n")
		if err := s.Transport.Close(); err != nil {
			log.Printf("transport close error: %s\n", err.Error())
		}
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Printf("gob decode error: %s\n", err.Error())
			}
			if err := s.handleMessage(rpc.From.String(), &msg); err != nil {
				log.Printf("handle message error: %s\n", err.Error())
			}

		case <-s.doneChan:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, v)
	case MessageGetFile:
		return s.handleMessageGetFile(from, v)
	}
	return nil
}

func (s *FileServer) handleMessageGetFile(from string, msg MessageGetFile) error {
	if !s.Storage.Has(msg.ID, msg.Key) {
		return fmt.Errorf("[%s] need to serve file (%s) but it does not exist on disk", s.Transport.Addr(), msg.Key) //nolint:err113
	}

	log.Printf("[%s] serving file (%s) over the network\n", s.Transport.Addr(), msg.Key)

	fileSize, r, err := s.Storage.Read(msg.ID, msg.Key)
	if err != nil {
		return err
	}

	if rc, ok := r.(io.ReadCloser); ok {
		defer func() { _ = rc.Close() }()
	}

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found in the peers map", from) //nolint:err113
	}

	// First send the "incomingStream" byte to the peer then
	// we can send the file size as an int64.
	_ = peer.Send([]byte{p2p.IncomingStream})
	if wErr := binary.Write(peer, binary.LittleEndian, fileSize); wErr != nil {
		return wErr
	}
	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}

	log.Printf("[%s] written (%d) bytes over the network to from %s\n", s.Transport.Addr(), n, from)

	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found in the peers map", from) //nolint:err113
	}

	n, err := s.Storage.Write(msg.ID, msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}
	log.Printf("[%s] written %d bytes to disk\n", s.Transport.Addr(), n)

	peer.CloseStream()

	return nil
}

func (s *FileServer) bootstrapNetwork() {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			log.Printf("[%s] attempting to connect with remote: %s\n", s.Transport.Addr(), addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Printf("dial error: %s\n", err.Error())
			}
		}(addr)
	}
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
