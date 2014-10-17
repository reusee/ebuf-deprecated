package ebuf

import (
	"github.com/reusee/rope"
)

type Buffer struct {
	rope          *rope.Rope
	Cursors       []Cursor
	adjustCursors []Cursor
	savedStates   []State
	savingState   bool
	redoStates    []State
	Scanners      []*Scanner
}

func New() *Buffer {
	return &Buffer{
		rope:        rope.NewFromBytes(nil),
		savingState: true,
		Cursors:     []Cursor{0},
	}
}

func (b *Buffer) Bytes() []byte {
	return b.rope.Bytes()
}

func (b *Buffer) Len() int {
	return b.rope.Len()
}

func (b *Buffer) SubBytes(cur Cursor, l int) []byte {
	if l == 1 {
		return []byte{b.rope.Index(cur.Int())}
	}
	return b.rope.Sub(cur.Int(), l)
}

func (b *Buffer) SetBytes(bs []byte) {
	b.saveState()
	b.rope = rope.NewFromBytes(bs)
	b.redoStates = nil
}

func (b *Buffer) Action(fn func()) {
	b.saveState()
	b.savingState = false
	fn()
	b.savingState = true
}

func (b *Buffer) Insert(cursor Cursor, bs []byte) {
	b.saveState()
	b.redoStates = nil
	b.rope = b.rope.Insert(cursor.Int(), bs)

	newCursors := make(map[Cursor]struct{})
	for _, c := range b.Cursors {
		if c >= cursor {
			newCursors[c.Move(len(bs))] = struct{}{}
		} else {
			newCursors[c] = struct{}{}
		}
	}
	for i, c := range b.adjustCursors {
		if c >= cursor {
			b.adjustCursors[i] = b.adjustCursors[i].Move(len(bs))
		}
	}
	cursors := make([]Cursor, 0, len(newCursors))
	for c, _ := range newCursors {
		cursors = append(cursors, c)
	}
	b.Cursors = cursors
}

func (b *Buffer) Delete(cursor Cursor, length int) {
	b.saveState()
	b.redoStates = nil
	b.rope = b.rope.Delete(cursor.Int(), length)

	newCursors := make(map[Cursor]struct{})
	for _, c := range b.Cursors {
		if c <= cursor {
			newCursors[c] = struct{}{}
		} else {
			newCursors[c.Move(-length)] = struct{}{}
		}
	}
	for i, c := range b.adjustCursors {
		if c > cursor {
			b.adjustCursors[i] = b.adjustCursors[i].Move(-length)
		}
	}
	cursors := make([]Cursor, 0, len(newCursors))
	for c, _ := range newCursors {
		cursors = append(cursors, c)
	}
	b.Cursors = cursors
}
