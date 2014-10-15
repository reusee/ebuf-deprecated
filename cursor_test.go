package ebuf

import (
	"bytes"
	"testing"
)

func TestCursorMove(t *testing.T) {
	if Cursor(3).Move(-5) != 0 {
		t.Fatal()
	}
}

func TestSetAndUnsetCursor(t *testing.T) {
	r := New()
	r.SetCursor(0)
	if len(r.Cursors) != 1 {
		t.Fatal()
	}
	r.SetCursor(1)
	r.SetCursor(1)
	if len(r.Cursors) != 2 {
		t.Fatal()
	}
	r.UnsetCursor(0)
	if len(r.Cursors) != 1 {
		t.Fatal()
	}
	r.UnsetCursor(0)
	if len(r.Cursors) != 1 {
		t.Fatal()
	}
}

func TestInsertAtCursors(t *testing.T) {
	r := New()
	r.SetBytes([]byte("foobarbaz"))
	for i := 0; i < r.Len(); i++ {
		r.SetCursor(i)
	}
	r.InsertAtCursors([]byte("|"))
	if !bytes.Equal(r.Bytes(), []byte("|f|o|o|b|a|r|b|a|z")) {
		t.Fatal()
	}

	r.Undo()
	if !bytes.Equal(r.Bytes(), []byte("foobarbaz")) {
		t.Fatal()
	}
	r.Redo()
	if !bytes.Equal(r.Bytes(), []byte("|f|o|o|b|a|r|b|a|z")) {
		t.Fatal()
	}
}
