package ebuf

import "testing"

func TestDisplayWidth(t *testing.T) {
	if DisplayWidth([]byte("foo")) != 3 {
		t.Fatal()
	}
	if DisplayWidth([]byte("玻璃")) != 4 {
		t.Fatal()
	}
	if DisplayWidth([]byte("じょじょ")) != 4 {
		t.Fatal()
	}
	if DisplayWidth([]byte("我能吞zuo下da玻si璃而不伤身体")) != 29 {
		t.Fatal()
	}
}
