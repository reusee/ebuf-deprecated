package ebuf

import "fmt"

var p = fmt.Printf

func reversedBytes(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, b := range bs {
		ret[len(bs)-i-1] = b
	}
	return ret
}
