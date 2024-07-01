package store

import (
	//nolint:gosec
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yigithankarabulut/distributed-file-storage/crypto"
)

const defaultRootFolderName = "store"

// Store is an interface that can be implemented to store.
type Store struct {
	// Root is the root directory of the store, containing all the folders/files of the system.
	Root              string
	PathTransformFunc PathTransformFunc
}

// Option is a functional option for configuring a Store.
type Option func(*Store)

// WithRoot is a functional option for setting the root of the Store.
func WithRoot(root string) Option {
	return func(s *Store) {
		s.Root = root
	}
}

// WithPathTransformFunc is a functional option for setting the path transform function of the Store.
func WithPathTransformFunc(f PathTransformFunc) Option {
	return func(s *Store) {
		s.PathTransformFunc = f
	}
}

// NewStore creates a new Store with the given options.
func NewStore(opts ...Option) *Store {
	s := &Store{
		Root:              defaultRootFolderName,
		PathTransformFunc: DefaultPathTransformFunc,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Clear clears the all folders/files in the storage.
func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

// Has checks if a key exists in the storage.
func (s *Store) Has(id, key string) bool {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	_, err := os.Stat(pathNameWithRoot)

	return !errors.Is(err, os.ErrNotExist)
}

// Delete deletes a key from the storage.
func (s *Store) Delete(id, key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk\n", pathKey.FullPath())
	}()

	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FirstPathName())
	return os.RemoveAll(pathNameWithRoot)
}

// Write writes a key to the storage.
func (s *Store) Write(id, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)
}

// WriteDecrypt writes a key to the storage with decryption. It uses the given encryption key to decrypt the data.
func (s *Store) WriteDecrypt(encryptKey []byte, id, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}

	n, err := crypto.CopyDecrypt(encryptKey, r, f)
	return int64(n), err
}

func (s *Store) openFileForWriting(id, key string) (*os.File, error) {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil { //nolint:gosec
		return nil, err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	return os.Create(fullPathWithRoot) //nolint:gosec
}

func (s *Store) writeStream(id, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}

	return io.Copy(f, r)
}

func (s *Store) Read(id, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)
}

func (s *Store) readStream(id, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	pathKeyWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	file, err := os.Open(pathKeyWithRoot) //nolint:gosec
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return fi.Size(), file, nil
}
