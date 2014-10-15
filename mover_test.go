package ebuf

import "testing"

func TestRuneMover(t *testing.T) {
	r := New()
	r.SetBytes([]byte("我能吞zuo下da玻si璃而不伤身体"))
	r.MoveCursors(RuneMover, 0)
	if r.Cursors[0] != 0 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover, 1)
	if r.Cursors[0] != 3 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover, 2)
	if r.Cursors[0] != 9 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover, 1)
	if r.Cursors[0] != 10 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover, 14)
	if r.Cursors[0].Int() != r.Len() {
		t.Fatal()
	}
	r.MoveCursors(RuneMover, 1) // overflow
	if r.Cursors[0].Int() != r.Len() {
		t.Fatal()
	}

	r.Cursors[0] = Cursor(r.Len() - 3)
	r.MoveCursors(RuneMover, 2) // overflow
	if r.Cursors[0].Int() != r.Len() {
		t.Fatal()
	}

	r.MoveCursors(RuneMover, -1)
	if r.Cursors[0].Int() != r.Len()-3 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover, -7)
	if r.Cursors[0].Int() != 20 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover, -12)
	if r.Cursors[0].Int() != 0 {
		t.Fatal()
	}
	r.MoveCursors(RuneMover, -20) // overflow
	if r.Cursors[0].Int() != 0 {
		t.Fatal()
	}

	r.MoveCursors(RuneMover, 1)
	r.MoveCursors(RuneMover, -2)
	if r.Cursors[0].Int() != 0 {
		p("%d\n", r.Cursors[0])
		t.Fatal()
	}
}
