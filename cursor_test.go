package ebuf

import "testing"

func TestCursorMove(t *testing.T) {
	if Cursor(3).Move(-5) != 0 {
		t.Fatal()
	}
}
