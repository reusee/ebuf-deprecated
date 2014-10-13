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

func (b *Buffer) getState() State {
	state := State{
		rope:    b.rope,
		cursors: make([]Cursor, len(b.cursors)),
	}
	copy(state.cursors[:], b.cursors[:])
	return state
}

func (b *Buffer) saveState() {
	if !b.savingState {
		return
	}
	b.savedStates = append(b.savedStates, b.getState())
}

func (b *Buffer) loadState(state State) {
	b.rope = state.rope
	b.cursors = make([]Cursor, len(state.cursors))
	copy(b.cursors, state.cursors)
}

func (b *Buffer) SetBytes(bs []byte) {
	b.saveState()
	b.rope = rope.NewFromBytes(bs)
	b.redoStates = nil
}

func (b *Buffer) Bytes() []byte {
	return b.rope.Bytes()
}

func (b *Buffer) Action(fn func()) {
	b.savingState = false
	fn()
	b.savingState = true
}
