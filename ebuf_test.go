package ebuf

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	New()
}

func TestBytes(t *testing.T) {
	b := New()
	bs := []byte("foobarbaz")
	b.SetBytes(bs)
	if !bytes.Equal(b.Bytes(), bs) {
		t.Fatal()
	}
}
