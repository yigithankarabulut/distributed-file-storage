package p2p

// Peer represents a connection to another node in the network.
type Peer interface {
	Close() error
}

// Transport represents a network transport that can listen for incoming
// connections and accept them, as well as send and receive messages.
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
