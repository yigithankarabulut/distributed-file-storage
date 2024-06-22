package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "some-key"
	pathKey := CASPathTransformFunc(key)
	expectedFileName := "9cea46b39bd44a1ef9f3e71bfe9e45c24d3300f6"
	expectedPathName := "9cea4/6b39b/d44a1/ef9f3/e71bf/e9e45/c24d3/300f6"

	if pathKey.PathName != expectedPathName {
		t.Errorf("Expected %s, got %s", expectedPathName, pathKey.PathName)
	}
	if pathKey.FileName != expectedFileName {
		t.Errorf("Expected %s, got %s", expectedFileName, pathKey.FileName)
	}
}

func TestDelete(t *testing.T) {
	s := NewStore(
		WithPathTransformFunc(CASPathTransformFunc),
	)
	key := "my-special-picture"
	data := []byte("some jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

	if _, err := s.Read(key); err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStore(t *testing.T) {
	s := NewStore(
		WithPathTransformFunc(CASPathTransformFunc),
	)
	key := "my-special-picture"
	data := []byte("some jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("Expected %s to exist", key)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	if !bytes.Equal(b, data) {
		t.Errorf("Expected %s, got %s", string(data), string(b))
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

	_ = s.Delete(key)
}
