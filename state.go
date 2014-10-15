package ebuf

import "github.com/reusee/rope"

type State struct {
	rope    *rope.Rope
	cursors []Cursor
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
