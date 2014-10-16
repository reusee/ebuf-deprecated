package ebuf

import (
	"bytes"
	"testing"
)

func TestRuneMover(t *testing.T) {
	r := New()
	r.SetBytes([]byte("我能吞zuo下da玻si璃而不伤身体"))
	r.MoveCursors(RuneMover(0))
	if r.Cursors[0] != 0 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover(1))
	if r.Cursors[0] != 3 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover(2))
	if r.Cursors[0] != 9 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover(1))
	if r.Cursors[0] != 10 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover(14))
	if r.Cursors[0].Int() != r.Len() {
		t.Fatal()
	}
	r.MoveCursors(RuneMover(1)) // overflow
	if r.Cursors[0].Int() != r.Len() {
		t.Fatal()
	}

	r.Cursors[0] = Cursor(r.Len() - 3)
	r.MoveCursors(RuneMover(2)) // overflow
	if r.Cursors[0].Int() != r.Len() {
		t.Fatal()
	}

	r.MoveCursors(RuneMover(-1))
	if r.Cursors[0].Int() != r.Len()-3 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover(-7))
	if r.Cursors[0].Int() != 20 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover(-12))
	if r.Cursors[0].Int() != 0 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover(-20)) // overflow
	if r.Cursors[0].Int() != 0 {
		t.Fatal()
	}

	r.MoveCursors(RuneMover(1))
	r.MoveCursors(RuneMover(-2))
	if r.Cursors[0].Int() != 0 {
		p("%d\n", r.Cursors[0])
		t.Fatal()
	}
}

func BenchmarkRuneMover(b *testing.B) {
	r := New()
	r.SetBytes(bytes.Repeat([]byte{'x'}, 5000000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.MoveCursors(RuneMover(1))
	}
}

func TestMatchMover(t *testing.T) {
	bs := bytes.Repeat([]byte("foo"), 512)
	r := New()
	r.SetBytes(bs)

	for i := 0; i < 511; i++ {
		r.MoveCursors(MatchMover([]byte("foo"), 2))
		if r.Cursors[0].Int() != (i+1)*3 {
			t.Fatal()
		}
	}
	r.MoveCursors(MatchMover([]byte("foo"), 2))
	if r.Cursors[0].Int() != 511*3 {
		t.Fatal()
	}

	r.Cursors[0] = Cursor(r.Len())
	for i := 0; i < 511; i++ {
		r.MoveCursors(MatchMover([]byte("foo"), -2))
		if r.Cursors[0].Int() != (511-i)*3 {
			t.Fatal()
		}
	}
	r.MoveCursors(MatchMover([]byte("foo"), -2))
	if r.Cursors[0].Int() != 3 {
		t.Fatal()
	}

	r.SetBytes([]byte("fofofoo"))
	r.Cursors[0] = 0
	r.MoveCursors(MatchMover([]byte("foo"), 1))
	if r.Cursors[0] != 4 {
		t.Fatal()
	}

	r.SetBytes([]byte("foofofo"))
	r.Cursors[0] = Cursor(r.Len())
	r.MoveCursors(MatchMover([]byte("foo"), -1))
	if r.Cursors[0] != 3 {
		t.Fatal()
	}
}
