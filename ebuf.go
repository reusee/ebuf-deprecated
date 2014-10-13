package ebuf

import (
	"github.com/reusee/rope"
)

type Buffer struct {
	rope        *rope.Rope
	cursors     []Cursor
	savedStates []State
	savingState bool
}

type Cursor int

type State struct {
	rope    *rope.Rope
	cursors []Cursor
}

func New() *Buffer {
	return &Buffer{
		rope:        rope.NewFromBytes(nil),
		savingState: true,
		cursors:     []Cursor{0},
	}
}

func (b *Buffer) saveState() {
	if !b.savingState {
		return
	}
	state := State{
		rope:    b.rope,
		cursors: make([]Cursor, len(b.cursors)),
	}
	copy(state.cursors[:], b.cursors[:])
	b.savedStates = append(b.savedStates, state)
}

func (b *Buffer) SetBytes(bs []byte) {
	b.saveState()
	b.rope = rope.NewFromBytes(bs)
}

func (b *Buffer) Bytes() []byte {
	return b.rope.Bytes()
}

func (b *Buffer) Action(fn func()) {
	b.savingState = false
	fn()
	b.savingState = true
}
