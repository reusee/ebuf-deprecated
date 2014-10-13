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
		rope:        rope.NewFromBytes(nil),
		savingState: true,
	}
}

func (b *Buffer) SetBytes(bs []byte) {
	if b.savingState {
		b.savedStates = append(b.savedStates, b.rope)
	}
	b.rope = rope.NewFromBytes(bs)
}

func (b *Buffer) Bytes() []byte {
	return b.rope.Bytes()
}
