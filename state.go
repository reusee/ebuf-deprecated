package ebuf

import "github.com/reusee/rope"

type State struct {
	rope     *rope.Rope
	cursors  []Cursor
	scanners []*Scanner
}

func (b *Buffer) getState() State {
	state := State{
		rope:     b.rope,
		cursors:  make([]Cursor, len(b.Cursors)),
		scanners: make([]*Scanner, len(b.Scanners)),
	}
	copy(state.cursors[:], b.Cursors[:])
	copy(state.scanners[:], b.Scanners[:])
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
	b.Cursors = make([]Cursor, len(state.cursors))
	copy(b.Cursors, state.cursors)
	b.Scanners = make([]*Scanner, len(state.scanners))
	copy(b.Scanners, state.scanners)
}
