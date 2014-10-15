package ebuf

import "sort"

type Cursor int

type Range struct {
	Begin, End Cursor
}

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

func (b *Buffer) MoveCursors(mover Mover, n int) {
	cursors := make(map[Cursor]struct{})
	for _, c := range b.Cursors {
		cursors[mover(b, c, n)] = struct{}{}
	}
	newCursors := make([]Cursor, 0, len(cursors))
	for c, _ := range cursors {
		newCursors = append(newCursors, c)
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

func (b *Buffer) DeleteAtCursors(mover Mover, n int) {
	ranges := []Range{}
	for _, c := range b.Cursors {
		stop := mover(b, c, n)
		if stop > c {
			ranges = append(ranges, Range{c, stop})
		} else {
			ranges = append(ranges, Range{stop, c})
		}
	}
	sort.Sort(RangesSorter(ranges))
	delRanges := []Range{}
	for i, r := range ranges {
		if i == 0 {
			delRanges = append(delRanges, r)
		} else {
			last := delRanges[len(delRanges)-1]
			if r.Begin >= last.Begin && r.Begin <= last.End { // overlapped
				if r.End > last.End {
					last.End = r.End
					delRanges[len(delRanges)-1] = last
				}
			} else {
				delRanges = append(delRanges, r)
			}
		}
	}
	b.Action(func() {
		adjustCursors := make([]Cursor, len(delRanges))
		delLens := make([]int, len(delRanges))
		for i, r := range delRanges {
			adjustCursors[i] = r.Begin
			delLens[i] = (r.End - r.Begin).Int()
		}
		b.adjustCursors = adjustCursors
		for i := 0; i < len(b.adjustCursors); i++ {
			b.Delete(b.adjustCursors[i], delLens[i])
		}
	})
}

type RangesSorter []Range

func (r RangesSorter) Len() int           { return len(r) }
func (r RangesSorter) Less(i, j int) bool { return r[i].Begin < r[j].Begin }
func (r RangesSorter) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
