package ebuf

import "github.com/reusee/rope"

type State struct {
	rope    *rope.Rope
	Cursors []Cursor
}

func (b *Buffer) getState() State {
	state := State{
		rope:    b.rope,
		Cursors: make([]Cursor, len(b.Cursors)),
	}
	copy(state.Cursors[:], b.Cursors[:])
	return state
}

func (b *Buffer) saveState() {
	if !b.savingState {
		return
	}
	b.savedStates = append(b.savedStates, b.getState())
}

func (b *Buffer) loadState(state State) {
	b.State = state
}
