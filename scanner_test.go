package ebuf

import "testing"

func TestScanner(t *testing.T) {
	r := New()
	r.SetBytes([]byte(`foobarbazbazfoobar`))

	scanner := NewScanner(map[string][]string{
		"foo": {"foo"},
		"bar": {"bar"},
		"baz": {"(b)(a)(z)"},
	})

	r.rope.IterRune(0, func(ru rune, l int) bool {
		scanner.FeedRune(ru, l)
		return true
	})

	expected := []Capture{
		{"foo", 0, 3},
		{"bar", 3, 6},
		{"baz", 6, 9},
		{"baz", 9, 12},
		{"foo", 12, 15},
		{"bar", 15, 18},
	}
	for i, c := range scanner.Captures {
		if c != expected[i] {
			t.Fatal()
		}
	}

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal()
			}
		}()
		r.SetScanner(map[string][]string{
			"foo": {"["},
		})
	}()

	scanner = NewScanner(map[string][]string{
		"foo": {"foo", "bar", "baz"},
	})
	r.rope.IterRune(0, func(ru rune, l int) bool {
		scanner.FeedRune(ru, l)
		return true
	})

	expected = []Capture{
		{"foo", 0, 3},
		{"foo", 3, 6},
		{"foo", 6, 9},
		{"foo", 9, 12},
		{"foo", 12, 15},
		{"foo", 15, 18},
	}
	for i, c := range scanner.Captures {
		if c != expected[i] {
			t.Fatal()
		}
	}

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal()
			}
		}()
		scanner = NewScanner(map[string][]string{
			"foo": {"^$"},
		})
		r.rope.IterRune(0, func(ru rune, l int) bool {
			scanner.FeedRune(ru, l)
			return true
		})
	}()

	scanner = NewScanner(map[string][]string{
		"foo": {"()"},
	})
	r.rope.IterRune(0, func(ru rune, l int) bool {
		scanner.FeedRune(ru, l)
		return true
	})
	if len(scanner.Captures) != 18 {
		t.Fatal()
	}

	scanner = NewScanner(map[string][]string{
		"foo": {".*"},
	})
	r.rope.IterRune(0, func(ru rune, l int) bool {
		scanner.FeedRune(ru, l)
		return true
	})
	if len(scanner.Captures) != 18 {
		t.Fatal()
	}

}

func BenchmarkScanner(b *testing.B) {
	scanner := NewScanner(map[string][]string{
		"zuo": {"zuo"},
		"da":  {"da"},
		"si":  {"si"},
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner.FeedRune('z', 1)
	}
}
