package ebuf

import "testing"

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
