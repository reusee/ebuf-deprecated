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

func BenchmarkInsertAtCursors(b *testing.B) {
	r := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.InsertAtCursors([]byte{'x'})
	}
}

func TestDeleteAtCursors(t *testing.T) {
	r := New()
	r.SetBytes([]byte("foobarbaz"))
	r.SetCursor(0)
	r.SetCursor(3)
	r.SetCursor(6)

	r.DeleteAtCursors(RuneMover(1))
	if !bytes.Equal(r.Bytes(), []byte("ooaraz")) {
		t.Fatal()
	}
	r.Undo()
	if !bytes.Equal(r.Bytes(), []byte("foobarbaz")) {
		t.Fatal()
	}
	r.Redo()

	r.DeleteAtCursors(RuneMover(1))
	if !bytes.Equal(r.Bytes(), []byte("orz")) {
		t.Fatal()
	}
	r.Undo()
	if !bytes.Equal(r.Bytes(), []byte("ooaraz")) {
		t.Fatal()
	}
	r.Redo()

	r.DeleteAtCursors(RuneMover(1))
	if !bytes.Equal(r.Bytes(), nil) {
		t.Fatal()
	}
	r.Undo()
	if !bytes.Equal(r.Bytes(), []byte("orz")) {
		t.Fatal()
	}
	r.Redo()

	if len(r.Cursors) != 1 && r.Cursors[0] == 1 {
		t.Fatal()
	}

	r = New()
	r.SetBytes([]byte("foobarbaz"))
	r.SetCursor(3)
	r.SetCursor(6)
	r.SetCursor(9)

	r.DeleteAtCursors(RuneMover(-1))
	if !bytes.Equal(r.Bytes(), []byte("fobaba")) {
		t.Fatal()
	}
	r.Undo()
	if !bytes.Equal(r.Bytes(), []byte("foobarbaz")) {
		t.Fatal()
	}
	r.Redo()

	r.DeleteAtCursors(RuneMover(-1))
	if !bytes.Equal(r.Bytes(), []byte("fbb")) {
		t.Fatal()
	}
	r.Undo()
	if !bytes.Equal(r.Bytes(), []byte("fobaba")) {
		t.Fatal()
	}
	r.Redo()

	r.DeleteAtCursors(RuneMover(-1))
	if !bytes.Equal(r.Bytes(), nil) {
		t.Fatal()
	}
	r.Undo()
	if !bytes.Equal(r.Bytes(), []byte("fbb")) {
		t.Fatal()
	}
	r.Redo()

	if len(r.Cursors) != 1 && r.Cursors[0] == 0 {
		t.Fatal()
	}
}

func BenchmarkDeleteAtCursors(b *testing.B) {
	r := New()
	r.SetBytes(bytes.Repeat([]byte{'x'}, 5000000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.DeleteAtCursors(RuneMover(1))
	}
}
