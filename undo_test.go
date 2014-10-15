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

func TestRedo(t *testing.T) {
	b := New()
	b.SetBytes([]byte("foo"))
	b.Undo()
	b.Redo()
	if !bytes.Equal(b.Bytes(), []byte("foo")) {
		t.Fatal()
	}
	b.Undo()
	b.SetBytes([]byte("bar")) // this should clear the redo states
	b.Redo()
	if !bytes.Equal(b.Bytes(), []byte("bar")) {
		t.Fatal()
	}
}
