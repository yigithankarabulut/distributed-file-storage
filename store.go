package main

import (
	"bytes"
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "store"

// CASPathTransformFunc is a function that transforms a key into a CAS path.
// CAS stands for Content Addressable Storage.
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key)) //nolint:gosec
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

// PathTransformFunc is a function that transforms a path.
type PathTransformFunc func(string) PathKey

// DefaultPathTransformFunc is the default path transform function.
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

// PathKey is a struct that contains the pathname and the original key.
type PathKey struct {
	PathName string
	FileName string
}

// FirstPathName returns the first path name of the PathKey.
func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// FullPath returns the filename of the PathKey.
func (p PathKey) FullPath() string {
	return p.PathName + "/" + p.FileName
}

// Store is an interface that can be implemented to store.
type Store struct {
	// Root is the root directory of the store, containing all the folders/files of the system.
	Root              string
	PathTransformFunc PathTransformFunc
}

// StoreOption is a functional option for configuring a Store.
type StoreOption func(*Store)

// WithRoot is a functional option for setting the root of the Store.
func WithRoot(root string) StoreOption {
	return func(s *Store) {
		s.Root = root
	}
}

// WithPathTransformFunc is a functional option for setting the path transform function of the Store.
func WithPathTransformFunc(f PathTransformFunc) StoreOption {
	return func(s *Store) {
		s.PathTransformFunc = f
	}
}

// NewStore creates a new Store with the given options.
func NewStore(opts ...StoreOption) *Store {
	s := &Store{
		PathTransformFunc: DefaultPathTransformFunc,
		Root:              defaultRootFolderName,
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
func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	_, err := os.Stat(pathNameWithRoot)

	return !errors.Is(err, os.ErrNotExist)
}

// Delete deletes a key from the storage.
func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk\n", pathKey.FullPath())
	}()

	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	return os.RemoveAll(pathNameWithRoot)
}

// Write writes a key to the storage.
func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	pathKeyWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	return os.Open(pathKeyWithRoot) //nolint:gosec
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil { //nolint:gosec
		return 0, err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	f, err := os.Create(fullPathWithRoot) //nolint:gosec
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}
