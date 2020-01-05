package eris

import (
	"fmt"
	"runtime"
	"strings"
)

// Stack is an array of stack frames stored in a human readable format.
type Stack []StackFrame

func (s *Stack) contains(frame StackFrame) bool {
	for _, f := range *s {
		if f == frame {
			return true
		}
	}
	return false
}

// todo: this needs some improvement and explanatory comments
func (s *Stack) combineStack(wStack []StackFrame) {
	if s.contains(wStack[0]) {
		return
	}
	rStack := *s
	if len(wStack) == 1 || len(rStack) == 1 {
		rStack = append(rStack, wStack[0])
	} else if len(wStack) > 1 {
		for i, f := range rStack {
			if f == wStack[1] {
				rStack = append(rStack[:i], append([]StackFrame{wStack[0]}, rStack[i:]...)...)
				break
			}
		}
	}
	*s = rStack
}

func (s *Stack) format(sep string) []string {
	var str []string
	for _, f := range *s {
		str = append(str, f.format(sep))
	}
	return str
}

// StackFrame stores a frame's runtime information in a human readable format.
type StackFrame struct {
	Name string
	File string
	Line int
}

func (f *StackFrame) format(sep string) string {
	return fmt.Sprintf("%v%v%v%v%v", f.Name, sep, f.File, sep, f.Line)
}

// callers returns a stack trace.
func callers() *stack {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0 : n-2]
	return &st
}

// frame is a single program counter of a stack frame.
type frame uintptr

func (f frame) get() StackFrame {
	pc := uintptr(f) - 1
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return StackFrame{
			Name: "unknown",
			File: "unknown",
		}
	}

	name := fn.Name()
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	file, line := fn.FileLine(pc)

	return StackFrame{
		Name: name,
		File: file,
		Line: line,
	}
}

// stack is an array of program counters.
type stack []uintptr

func (s *stack) get() []StackFrame {
	var sFrames []StackFrame
	for _, f := range *s {
		frame := frame(f)
		sFrame := frame.get()
		sFrames = append(sFrames, sFrame)
	}
	return sFrames
}
