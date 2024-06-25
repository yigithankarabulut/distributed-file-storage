package store

import (
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
	"strings"
)

// PathTransformFunc is a function that transforms a path.
type PathTransformFunc func(string) PathKey

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

// DefaultPathTransformFunc is the default path transform function.
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}
