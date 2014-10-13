package ebuf

func (b *Buffer) Undo() {
	if len(b.savedStates) == 0 {
		return
	}
	currentState := b.getState()
	b.redoStates = append(b.redoStates, currentState)
	state := b.savedStates[len(b.savedStates)-1]
	b.savedStates = b.savedStates[:len(b.savedStates)-1]
	b.loadState(state)
}

func (b *Buffer) Redo() {
	if len(b.redoStates) == 0 {
		return
	}
	b.saveState()
	state := b.redoStates[len(b.redoStates)-1]
	b.redoStates = b.redoStates[:len(b.redoStates)-1]
	b.loadState(state)
}
