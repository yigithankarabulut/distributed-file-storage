package fileserver

// Message is a struct that contains the payload of the message.
type Message struct {
	Payload any
}

// MessageStoreFile is a struct that contains the key and the size of the file.
type MessageStoreFile struct {
	Key  string
	Size int64
}

// MessageGetFile is a struct that contains the key of the file.
type MessageGetFile struct {
	Key string
}
