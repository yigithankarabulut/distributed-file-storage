package store

import "strings"

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
