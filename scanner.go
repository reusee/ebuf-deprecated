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
	freeThreads  []*_Thread
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
				var newThread *_Thread
				if len(s.freeThreads) > 0 {
					newThread = s.freeThreads[len(s.freeThreads)-1]
					s.freeThreads = s.freeThreads[:len(s.freeThreads)-1]
				} else {
					newThread = &_Thread{
						captures: make([]int, s.maxCapture+1),
					}
				}
				copy(newThread.captures, thread.captures)
				newThread.pc = inst.Arg
				threads = append(threads, newThread)
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
				if len(blockingThreads) > 0 {
					blockingThreads = blockingThreads[0:0]
				}
				s.freeThreads = append(s.freeThreads, thread)
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
						s.freeThreads = append(s.freeThreads, thread)
						break runLoop
					}
				}
			}
		}
	}

	if len(blockingThreads) == 0 {
		s.freeThreads = append(s.freeThreads, threads...) // free all threads
		if len(threads) > 0 {
			threads = threads[0:0]
		}
		var newThread *_Thread
		if len(s.freeThreads) > 0 {
			newThread = s.freeThreads[len(s.freeThreads)-1]
			s.freeThreads = s.freeThreads[:len(s.freeThreads)-1]
		} else {
			newThread = &_Thread{
				captures: make([]int, s.maxCapture+1),
			}
		}
		copy(newThread.captures, s.zeroCapture)
		newThread.pc = uint32(s.program.Start)
		blockingThreads = append(blockingThreads, newThread)
	}
	s.tmpThreads = threads
	s.threads = blockingThreads

	s.pos += l
}
