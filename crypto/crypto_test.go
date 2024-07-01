package crypto

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncrypt(t *testing.T) {
	payLoad := []byte("Foo not Bar")
	src := bytes.NewReader(payLoad)
	dst := new(bytes.Buffer)
	key, _ := NewEncryptionKey()
	if _, err := CopyEncrypt(key, src, dst); err != nil {
		t.Error(err)
	}

	fmt.Println(dst.String())

	out := new(bytes.Buffer)
	nw, err := CopyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}

	if nw != 16+len(payLoad) {
		t.Errorf("decrypted data size is not equal to original data size, got: %d, want: %d", nw, 16+len(payLoad))
	}

	if !bytes.Equal(out.Bytes(), payLoad) {
		t.Errorf("decrypted data is not equal to original data, got: %s, want: %s", out.String(), payLoad)
	}

	fmt.Println(out.String())
}
