package ebuf

import (
	"fmt"
	"regexp/syntax"
	"strings"
)

type Scanner struct {
	program      *syntax.Prog
	captureNames []string
	threads      []*_Thread
	tmpThreads   []*_Thread
	pos          int
	Captures     []Capture
}

type _Thread struct {
	pc       uint32
	captures map[uint32]int
}

type Capture struct {
	Name  string
	Begin Cursor
	End   Cursor
}

func (b *Buffer) SetScanner(rules map[string][]string) {
	b.Scanners = []*Scanner{
		NewScanner(rules),
	}
}

func NewScanner(rules map[string][]string) *Scanner {
	exprs := []string{}
	for name, subs := range rules {
		exprs = append(exprs, fmt.Sprintf(`(?P<%s>%s)`, name,
			strings.Join(subs, "|")))
	}
	expr := strings.Join(exprs, "|")

	re, err := syntax.Parse(expr, syntax.Perl)
	if err != nil {
		panic(fmt.Errorf("rule syntax error %v", err))
	}
	capNames := re.CapNames()
	re = re.Simplify()
	program, _ := syntax.Compile(re)

	return &Scanner{
		program:      program,
		captureNames: capNames,
		threads: []*_Thread{
			&_Thread{
				pc:       uint32(program.Start),
				captures: make(map[uint32]int),
			},
		},
	}
}

func (s *Scanner) FeedRune(r rune, l int) {
	threads := s.threads
	blockingThreads := s.tmpThreads
loop:
	for len(threads) > 0 {
		thread := threads[len(threads)-1]
		threads = threads[:len(threads)-1]
		runeConsumed := false
		pc := thread.pc
		inst := s.program.Inst[pc]
	runLoop:
		for {
			switch inst.Op {
			case syntax.InstAlt, syntax.InstAltMatch:
				// new thread
				captures := make(map[uint32]int, len(thread.captures))
				for nameIndex, pos := range thread.captures {
					captures[nameIndex] = pos
				}
				threads = append(threads, &_Thread{
					pc:       inst.Arg,
					captures: captures,
				})
				pc = inst.Out
				inst = s.program.Inst[pc]
			case syntax.InstCapture:
				nameIndex := inst.Arg / 2
				name := s.captureNames[nameIndex]
				if name != "" { // skip nameless groups
					if pos, ok := thread.captures[nameIndex]; ok { // end of named group
						s.Captures = append(s.Captures, Capture{
							Name:  s.captureNames[nameIndex],
							Begin: Cursor(pos),
							End:   Cursor(s.pos + l),
						})
						delete(thread.captures, nameIndex)
					} else {
						thread.captures[nameIndex] = s.pos
					}
				}
				pc = inst.Out
				inst = s.program.Inst[pc]
			case syntax.InstEmptyWidth:
				panic("empty string pattern is not supported")
			case syntax.InstMatch, syntax.InstFail: // clear all threads, restart
				blockingThreads = nil
				break loop
			case syntax.InstNop:
				pc = inst.Out
				inst = s.program.Inst[pc]
			case syntax.InstRune1, syntax.InstRune, syntax.InstRuneAny, syntax.InstRuneAnyNotNL:
				if runeConsumed { // thread blocks
					thread.pc = pc
					blockingThreads = append(blockingThreads, thread)
					break runLoop
				} else { // consume rune
					if inst.Op == syntax.InstRune1 && r == inst.Rune[0] ||
						inst.Op == syntax.InstRune && inst.MatchRune(r) ||
						inst.Op == syntax.InstRuneAny ||
						inst.Op == syntax.InstRuneAnyNotNL && r != '\n' { // rune matchs
						runeConsumed = true
						pc = inst.Out
						inst = s.program.Inst[pc]
					} else { // thread dies
						break runLoop
					}
				}
			}
		}
	}
	s.tmpThreads = threads

	if len(blockingThreads) > 0 {
		s.threads = blockingThreads
	} else { // restart
		s.threads = []*_Thread{
			{
				pc:       uint32(s.program.Start),
				captures: make(map[uint32]int),
			},
		}
	}

	s.pos += l
}
