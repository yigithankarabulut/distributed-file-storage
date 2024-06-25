package p2p

import (
	"encoding/gob"
	"io"
)

// Decoder is an interface that can be implemented to decode
// a message from a reader into an RPC message.
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

// GOBDecoder is a decoder that uses the gob package to decode
// messages from a reader into an RPC message.
type GOBDecoder struct{}

// Decode decodes a message from a reader into an RPC message.
// Implements the Decoder interface.
func (d *GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

// DefaultDecoder is a decoder that reads from a reader into an RPC message.
// It reads up to 1028 bytes from the reader and sets the payload of the RPC message.
type DefaultDecoder struct{}

// Decode decodes a message from a reader into an RPC message.
// Implements the Decoder interface.
func (d DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	peakBuf := make([]byte, 1)
	if _, err := r.Read(peakBuf); err != nil {
		return nil //nolint:nilerr
	}

	// In case of a stream we are not decoding what is being sent over the network
	// we are just setting stream true, so we can handle that in our logic.
	stream := peakBuf[0] == IncomingStream
	if stream {
		msg.Stream = true
		return nil
	}

	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
