package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5" //nolint:gosec
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// GenerateID generates a new unique ID.
func GenerateID() string {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}

// HashKey generates a hash key for the given key.
func HashKey(key string) string {
	hash := md5.Sum([]byte(key)) //nolint:gosec
	return hex.EncodeToString(hash[:])
}

// NewEncryptionKey generates a new encryption key.
func NewEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// CopyDecrypt reads from src, decrypts the data using the given key and writes to dst.
func CopyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// Read the IV from given io.Reader which, in our case should be the
	// block.BlockSize() bytes.
	iv := make([]byte, block.BlockSize())
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

// CopyEncrypt reads from src, encrypts the data using the given key and writes to dst.
func CopyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer) (int, error) {
	var (
		buf = make([]byte, 1024*32)
		nw  = blockSize
	)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			nn, err2 := dst.Write(buf[:n])
			if err2 != nil {
				return 0, err
			}
			nw += nn
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return nw, nil
}
