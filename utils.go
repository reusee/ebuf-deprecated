package ebuf

import (
	"fmt"
	"unicode/utf8"
)

var p = fmt.Printf

func reversedBytes(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, b := range bs {
		ret[len(bs)-i-1] = b
	}
	return ret
}

func DisplayWidth(bs []byte) (ret int) {
	for {
		r, n := utf8.DecodeRune(bs)
		if r == utf8.RuneError {
			break
		}
		bs = bs[n:]
		switch {
		case r >= 0x4e00 && r <= 0x9fff,
			r >= 0x3400 && r <= 0x4dbf,
			r >= 0xf900 && r <= 0xfaff,
			r >= 0x20000 && r <= 0x2ffff,
			r >= 0x30000 && r <= 0x3ffff:
			ret += 2
		default:
			ret += 1
		}
	}
	return
}
