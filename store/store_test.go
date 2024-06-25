package store

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "some-key"
	pathKey := CASPathTransformFunc(key)
	expectedFileName := "9cea46b39bd44a1ef9f3e71bfe9e45c24d3300f6"
	expectedPathName := "9cea4/6b39b/d44a1/ef9f3/e71bf/e9e45/c24d3/300f6"

	if pathKey.PathName != expectedPathName {
		t.Errorf("expected %s, got %s", expectedPathName, pathKey.PathName)
	}
	if pathKey.FileName != expectedFileName {
		t.Errorf("expected %s, got %s", expectedFileName, pathKey.FileName)
	}
}

func TestDelete(t *testing.T) {
	s := NewStore(
		WithPathTransformFunc(CASPathTransformFunc),
	)
	key := "my-special-picture"
	data := []byte("some jpg bytes")

	if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

	if _, _, err := s.Read(key); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(s, t)

	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("foo-%d", i)
		data := []byte("some jpg bytes")

		if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expected %s to exist", key)
		}

		_, r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if !bytes.Equal(b, data) {
			t.Errorf("expected %s, got %s", string(data), string(b))
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected to not have key %s", key)
		}
	}
}

func newStore() *Store {
	return NewStore(
		WithPathTransformFunc(CASPathTransformFunc),
	)
}

func teardown(s *Store, t *testing.T) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
