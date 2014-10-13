package ebuf

import (
	"bytes"
	"testing"
)

func TestUndo(t *testing.T) {
	b := New()
	b.SetBytes([]byte("foo"))
	b.Undo()
	if !bytes.Equal(b.Bytes(), []byte{}) {
		t.Fatal()
	}
	b.Redo()
	if !bytes.Equal(b.Bytes(), []byte("foo")) {
		t.Fatal()
	}
	b.Undo()
	if !bytes.Equal(b.Bytes(), []byte{}) {
		t.Fatal()
	}
	b.Undo()
	if !bytes.Equal(b.Bytes(), []byte{}) {
		t.Fatal()
	}
	b.Redo()
	if !bytes.Equal(b.Bytes(), []byte("foo")) {
		t.Fatal()
	}
	b.Redo()
	if !bytes.Equal(b.Bytes(), []byte("foo")) {
		t.Fatal()
	}
}
