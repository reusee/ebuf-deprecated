package ebuf

type Object interface {
	Move(*Buffer, Cursor) Cursor
}

//TODO
type RuneObject struct{}

//TODO
type LineObject struct{}

//TODO
type RegexSeparatedObject struct{}
