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

func TestAction(t *testing.T) {
	b := New()
	b.Action(func() {
		b.SetBytes([]byte("foobarbaz"))
		b.SetBytes([]byte("foobarbaz"))
		b.SetBytes([]byte("foobarbaz"))
		b.SetBytes([]byte("foobarbaz"))
	})
	if len(b.savedStates) > 1 {
		t.Fatal()
	}
}

func TestInsert(t *testing.T) {
	b := New()
	b.Insert(b.Cursors[0], []byte("foo"))
	if !bytes.Equal(b.Bytes(), []byte("foo")) {
		t.Fatal()
	}
	b.Insert(b.Cursors[0], []byte("bar"))
	if !bytes.Equal(b.Bytes(), []byte("foobar")) {
		t.Fatal()
	}
	b.Insert(Cursor(0), []byte("baz"))
	if !bytes.Equal(b.Bytes(), []byte("bazfoobar")) {
		t.Fatal()
	}
	if b.Cursors[0] != 9 {
		t.Fatal()
	}
	b.Cursors[0] = b.Cursors[0].Move(-3)
	b.Insert(Cursor(b.Len()), []byte("qux"))
	if !bytes.Equal(b.Bytes(), []byte("bazfoobarqux")) {
		t.Fatal()
	}
	if b.Cursors[0] != 6 {
		t.Fatal()
	}
}

func TestDelete(t *testing.T) {
	b := New()
	b.SetBytes([]byte("foobarbaz"))
	b.Delete(Cursor(0), 1)
	if !bytes.Equal(b.Bytes(), []byte("oobarbaz")) {
		t.Fatal()
	}
	if b.Cursors[0] != 0 {
		t.Fatal()
	}
	b.Delete(Cursor(0), 1)
	if !bytes.Equal(b.Bytes(), []byte("obarbaz")) {
		t.Fatal()
	}
	if b.Cursors[0] != 0 {
		t.Fatal()
	}

	b.Cursors[0] = b.Cursors[0].Move(b.Len())
	b.Delete(Cursor(0), 1)
	if !bytes.Equal(b.Bytes(), []byte("barbaz")) {
		t.Fatal()
	}
	if b.Cursors[0].Int() != b.Len() {
		t.Fatal()
	}
}

func BenchmarkRuneMover(b *testing.B) {
	r := New()
	r.SetBytes([]byte("我能吞zuo下da玻si璃而不伤身体"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Cursors[0] = 0
		r.MoveCursors(RuneMover, 1)
		r.MoveCursors(RuneMover, 1)
		r.MoveCursors(RuneMover, 1)
		r.MoveCursors(RuneMover, 1)
		r.MoveCursors(RuneMover, 1)
		r.MoveCursors(RuneMover, -1)
		r.MoveCursors(RuneMover, -1)
		r.MoveCursors(RuneMover, -1)
		r.MoveCursors(RuneMover, -1)
		r.MoveCursors(RuneMover, -1)
	}
}
