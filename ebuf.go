package ebuf

import (
	"github.com/reusee/rope"
)

type Buffer struct {
	rope        *rope.Rope
	savedStates []*rope.Rope
	savingState bool
}

func New() *Buffer {
	return &Buffer{
		rope: rope.NewFromBytes(nil),
	}
}
