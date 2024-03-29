package ebuf

import (
	"unicode/utf8"
)

type Mover func(*Buffer, Cursor) Cursor

func RuneMover(n int) Mover {
	return func(buf *Buffer, cur Cursor) Cursor {
		if n > 0 {
			bs := buf.SubBytes(cur, n*4)
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
			bs := buf.SubBytes(start, l.Int())
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
}

func MatchMover(bs []byte, n int, passthrough bool) Mover {
	return func(buf *Buffer, cur Cursor) Cursor {
		if n > 0 {
			bsIndex := 0
			offset := 0
			start := 0
			matched := false
			buf.rope.Iter(cur.Int(), func(slice []byte) bool {
				for _, b := range slice {
					if b == bs[bsIndex] {
						if bsIndex == 0 {
							start = offset
						}
						bsIndex++
					} else {
						bsIndex = 0
					}
					if bsIndex == len(bs) { // matched
						matched = true
						n--
						if n == 0 {
							return false
						} else {
							bsIndex = 0
						}
					}
					offset++
				}
				return true
			})
			if matched && passthrough {
				start += len(bs)
			}
			return cur.Move(start)
		} else {
			bs = reversedBytes(bs)
			bsIndex := 0
			offset := 0
			start := 0
			buf.rope.IterBackward(cur.Int(), func(slice []byte) bool {
				for _, b := range slice {
					if b == bs[bsIndex] {
						if bsIndex == 0 {
							start = offset
						}
						bsIndex++
					} else {
						bsIndex = 0
					}
					if bsIndex == len(bs) {
						n++
						if n == 0 {
							return false
						} else {
							bsIndex = 0
						}
					}
					offset++
				}
				return true
			})
			if passthrough {
				start += len(bs)
			}
			return cur.Move(-start)
		}
	}
}

func DisplayWidthMover(n int) Mover {
	return func(buf *Buffer, cur Cursor) Cursor {
		if n > 0 {
			bs := buf.SubBytes(cur, n*4)
			offset := 0
			for {
				r, l := utf8.DecodeRune(bs)
				if r == utf8.RuneError {
					break
				}
				bs = bs[l:]
				n -= RuneDisplayWidth(r)
				if n < 0 {
					break
				}
				offset += l
			}
			return cur.Move(offset)
		} else {
			start := cur.Move(n * 4)
			bs := buf.SubBytes(start, (cur - start).Int())
			offset := 0
			for {
				r, l := utf8.DecodeLastRune(bs)
				if r == utf8.RuneError {
					break
				}
				bs = bs[:len(bs)-l]
				n += RuneDisplayWidth(r)
				if n > 0 {
					break
				}
				offset += l
			}
			return cur.Move(-offset)
		}
	}
}
