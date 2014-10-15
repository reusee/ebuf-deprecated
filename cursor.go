package ebuf

type Cursor int

func (c Cursor) Int() int {
	return int(c)
}

func (c Cursor) Move(n int) Cursor {
	n = c.Int() + n
	if n < 0 {
		n = 0
	}
	return Cursor(n)
}

func (b *Buffer) SetCursor(offset int) {
	exists := false
	cursor := Cursor(offset)
	for _, c := range b.Cursors {
		if c == cursor {
			exists = true
		}
	}
	if !exists {
		b.Cursors = append(b.Cursors, cursor)
	}
}

func (b *Buffer) UnsetCursor(offset int) {
	newCursors := make([]Cursor, 0, len(b.Cursors)-1)
	cursor := Cursor(offset)
	for _, c := range b.Cursors {
		if c != cursor {
			newCursors = append(newCursors, c)
		}
	}
	b.Cursors = newCursors
}

func (b *Buffer) InsertAtCursors(bs []byte) {
	b.Action(func() {
		adjustCursors := make([]Cursor, len(b.Cursors))
		copy(adjustCursors[:], b.Cursors[:])
		b.adjustCursors = adjustCursors
		for i := 0; i < len(b.adjustCursors); i++ {
			b.Insert(b.adjustCursors[i], bs)
		}
	})
}

func (b *Buffer) DeleteAtCursors(obj Object, count int) {
	//TODO
}

func (b *Buffer) MoveCursors(obj Object, count int) {
	//TODO
}
