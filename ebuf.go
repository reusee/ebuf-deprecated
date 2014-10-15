package ebuf

import (
	"github.com/reusee/rope"
)

type Buffer struct {
	rope        *rope.Rope
	cursors     []Cursor
	savedStates []State
	savingState bool
	redoStates  []State
}

func New() *Buffer {
	return &Buffer{
		rope:        rope.NewFromBytes(nil),
		savingState: true,
		cursors:     []Cursor{0},
	}
}

func (b *Buffer) Bytes() []byte {
	return b.rope.Bytes()
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

func (b *Buffer) Insert(index int, bs []byte) {
	b.saveState()
	b.redoStates = nil
	//TODO
}

func (b *Buffer) Delete(index, length int) {
	b.saveState()
	b.redoStates = nil
	//TODO
}
