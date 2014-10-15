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

func TestSubBytes(t *testing.T) {
	bs := bytes.Repeat([]byte("foobarbaz"), 1)
	r := New()
	r.SetBytes(bs)
	for i := 0; i < len(bs); i++ {
		for j := 0; j < len(bs)-i; j++ {
			if !bytes.Equal(r.SubBytes(Cursor(i), j), bs[i:i+j]) {
				t.Fatal()
			}
		}
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
