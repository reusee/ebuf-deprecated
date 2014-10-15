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

func (b *Buffer) InsertAtCursors(bs []byte) {
	//TODO
}

func (b *Buffer) DeleteAtCursors(obj Object, count int) {
	//TODO
}

func (b *Buffer) MoveCursors(obj Object, count int) {
	//TODO
}
