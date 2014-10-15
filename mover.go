package ebuf

import (
	"unicode/utf8"
)

type Mover func(*Buffer, Cursor, int) Cursor

func RuneMover(buf *Buffer, cur Cursor, n int) Cursor {
	if n > 0 {
		bs := buf.rope.Sub(cur.Int(), n*4)
		offset := 0
		for i := 0; i < n; i++ {
			r, l := utf8.DecodeRune(bs)
			if r == utf8.RuneError {
				break
			}
			offset += l
			bs = bs[l:]
		}
		return cur.Move(offset)
	} else {
		start := cur.Move(n * 4)
		l := cur - start
		bs := buf.rope.Sub(start.Int(), l.Int())
		offset := 0
		for i := 0; i < -n; i++ {
			r, l := utf8.DecodeLastRune(bs)
			if r == utf8.RuneError {
				break
			}
			offset += l
			bs = bs[:len(bs)-l]
		}
		return cur.Move(-offset)
	}
}
