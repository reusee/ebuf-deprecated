package ebuf

type Cursor int

func (c Cursor) Int() int {
	return int(c)
}

func (c Cursor) Move(n int) Cursor {
	return c + Cursor(n)
}
