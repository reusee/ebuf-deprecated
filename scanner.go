package ebuf

import (
	"fmt"
	"regexp/syntax"
	"strings"
)

type Scanner struct {
	program      *syntax.Prog
	captureNames []string
	maxCapture   int
	zeroCapture  []int
	capturesPool [][]int
	threads      []*_Thread
	tmpThreads   []*_Thread
	pos          int
	Captures     []Capture
}

type _Thread struct {
	pc       uint32
	captures []int
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
	maxCapture := re.MaxCap()
	re = re.Simplify()
	program, _ := syntax.Compile(re)

	return &Scanner{
		program:      program,
		captureNames: capNames,
		maxCapture:   maxCapture,
		zeroCapture:  make([]int, maxCapture+1),
		threads: []*_Thread{
			&_Thread{
				pc:       uint32(program.Start),
				captures: make([]int, maxCapture+1),
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
				var captures []int
				if len(s.capturesPool) > 0 {
					captures = s.capturesPool[len(s.capturesPool)-1]
					s.capturesPool = s.capturesPool[:len(s.capturesPool)-1]
				} else {
					captures = make([]int, s.maxCapture+1)
				}
				copy(captures, thread.captures)
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
					if pos := thread.captures[nameIndex]; pos > 0 { // end of group
						s.Captures = append(s.Captures, Capture{
							Name:  s.captureNames[nameIndex],
							Begin: Cursor(pos - 1),
							End:   Cursor(s.pos + l),
						})
						thread.captures[nameIndex] = 0
					} else {
						thread.captures[nameIndex] = s.pos + 1 // 0 means not capturing
					}
				}
				pc = inst.Out
				inst = s.program.Inst[pc]
			case syntax.InstEmptyWidth:
				panic("empty string pattern is not supported")
			case syntax.InstMatch, syntax.InstFail: // clear all threads, restart
				blockingThreads = nil
				s.capturesPool = append(s.capturesPool, thread.captures)
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
						s.capturesPool = append(s.capturesPool, thread.captures)
						break runLoop
					}
				}
			}
		}
	}
	s.tmpThreads = threads

	if len(blockingThreads) == 0 {
		var captures []int
		if len(s.capturesPool) > 0 {
			captures = s.capturesPool[len(s.capturesPool)-1]
			s.capturesPool = s.capturesPool[:len(s.capturesPool)-1]
			copy(captures, s.zeroCapture)
		} else {
			captures = make([]int, s.maxCapture+1)
		}
		blockingThreads = append(blockingThreads, &_Thread{
			pc:       uint32(s.program.Start),
			captures: captures,
		})
	}
	s.threads = blockingThreads

	s.pos += l
}
