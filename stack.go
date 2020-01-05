package eris

import (
	"fmt"
	"runtime"
	"strings"
)

// Stack is an array of stack frames stored in a human readable format.
type Stack []StackFrame

// contains checks if a stack trace contains a specific frame.
func (s *Stack) contains(frame StackFrame) bool {
	for _, f := range *s {
		if f == frame {
			return true
		}
	}
	return false
}

// insertFrames inserts a wrap error frame into the correct place of the root error stack trace.
func (s *Stack) insertFrame(wFrames []StackFrame) {
	if len(wFrames) == 0 || s.contains(wFrames[0]) {
		return
	} else if len(wFrames) == 1 {
		// append the frame to the end if there's only one
		*s = append(*s, wFrames[0])
	}

	rStack := *s
	for i, f := range rStack {
		if f == wFrames[1] {
			// insert the first frame into the stack if the second frame is found
			rStack = append(rStack[:i], append([]StackFrame{wFrames[0]}, rStack[i:]...)...)
			break
		}
	}
	*s = rStack
}

// format returns an array of formatted stack frames.
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

// format returns a formatted stack frame.
func (f *StackFrame) format(sep string) string {
	return fmt.Sprintf("%v%v%v%v%v", f.Name, sep, f.File, sep, f.Line)
}

// callers returns a stack trace.
func callers(skip int) *stack {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0 : n-2]
	return &st
}

// frame is a single program counter of a stack frame.
type frame uintptr

// get returns a human readable stack frame.
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

// get returns a human readable stack trace.
func (s *stack) get() []StackFrame {
	var sFrames []StackFrame
	for _, f := range *s {
		frame := frame(f)
		sFrame := frame.get()
		sFrames = append(sFrames, sFrame)
	}
	return sFrames
}
