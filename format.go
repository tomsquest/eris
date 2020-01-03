package eris

import (
	"fmt"
)

// Format defines an error output format to be used with the default formatter.
type Format struct {
	WithTrace bool   // Flag that enables stack trace output.
	Msg       string // Separator between error messages and stack frame data.
	TBeg      string // Separator at the beginning of each stack frame.
	TSep      string // Separator between elements of each stack frame.
	Sep       string // Separator between each error in the chain.
}

// NewDefaultFormat conveniently returns a basic format for the default string formatter.
func NewDefaultFormat(withTrace bool) Format {
	stringFmt := Format{
		WithTrace: withTrace,
		Sep:       ": ",
	}
	if withTrace {
		stringFmt.Msg = "\n"
		stringFmt.TBeg = "\t"
		stringFmt.TSep = ": "
		stringFmt.Sep = "\n"
	}
	return stringFmt
}

// UnpackedError represents complete information about an error.
//
// This type can be used for custom error logging and parsing. Use `eris.Unpack` to build an UnpackedError
// from any error type. The ErrChain and ErrRoot fields correspond to `wrapError` and `rootError` types,
// respectively. If any other error type is unpacked, it will appear in the ExternalErr field.
type UnpackedError struct {
	ErrRoot     ErrRoot
	ErrChain    []ErrLink
	ExternalErr string
}

// Unpack returns UnpackedError type for a given golang error type.
func Unpack(err error) UnpackedError {
	e := UnpackedError{}
	switch err.(type) {
	case nil:
		return e
	case *rootError:
		e.unpackRootErr(err.(*rootError))
	case *wrapError:
		e.unpackWrapErr(err.(*wrapError))
	default:
		e.ExternalErr = err.Error()
	}
	return e
}

// ToString returns a default formatted string for a given eris error.
func (upErr *UnpackedError) ToString(format Format) string {
	var str string
	if len(upErr.ErrRoot.Msg) != 0 || len(upErr.ErrRoot.Stack) != 0 {
		str += upErr.ErrRoot.formatStr(format)
	}
	for _, eLink := range upErr.ErrChain {
		if !format.WithTrace {
			str += format.Sep
		}
		str += eLink.formatStr(format)
	}
	if upErr.ExternalErr != "" {
		str += fmt.Sprint(upErr.ExternalErr)
	}
	return str
}

// ToJSON returns a JSON formatted map for a given eris error.
func (upErr *UnpackedError) ToJSON(format Format) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	if len(upErr.ErrRoot.Msg) != 0 || len(upErr.ErrRoot.Stack) != 0 {
		jsonMap["root"] = upErr.ErrRoot.formatJSON(format)
	}
	if len(upErr.ErrChain) != 0 {
		var wrapArr []map[string]interface{}
		for _, eLink := range upErr.ErrChain {
			wrapMap := eLink.formatJSON(format)
			wrapArr = append(wrapArr, wrapMap)
		}
		jsonMap["wrap"] = wrapArr
	}
	if upErr.ExternalErr != "" {
		jsonMap["external"] = fmt.Sprint(upErr.ExternalErr)
	}
	return jsonMap
}

// unpackRootErr unpacks a rootError's message and stack trace.
// it also appends any additional wrapError frames to the stack.
func (upErr *UnpackedError) unpackRootErr(err *rootError) {
	stack := err.stack.get()
	for i := len(upErr.ErrChain) - 1; i >= 0; i-- {
		if !stackContains(stack, upErr.ErrChain[i].Frame) {
			stack = append(stack, upErr.ErrChain[i].Frame)
		}
	}
	upErr.ErrRoot.Msg = err.msg
	upErr.ErrRoot.Stack = stack
}

// unpackWrapErr unpacks a wrapError until it hits a rootError.
func (upErr *UnpackedError) unpackWrapErr(err *wrapError) {
	// Prepend each link so they'll appear in the same order as the stack.
	link := ErrLink{
		Msg:   err.msg,
		Frame: *err.frame.get(),
	}
	upErr.ErrChain = append([]ErrLink{link}, upErr.ErrChain...)

	nextErr := err.Unwrap()
	switch nextErr.(type) {
	case nil:
		return
	case *rootError:
		upErr.unpackRootErr(nextErr.(*rootError))
	case *wrapError:
		upErr.unpackWrapErr(nextErr.(*wrapError))
	default:
		upErr.ExternalErr = err.Error()
	}
}

// ErrRoot represents an error stack and the accompanying message.
type ErrRoot struct {
	Msg   string
	Stack []StackFrame
}

func (err *ErrRoot) formatStr(format Format) string {
	if err == nil {
		return ""
	}
	str := err.Msg
	str += format.Msg
	if format.WithTrace {
		stackArr := formatStackFrames(err.Stack, format.TSep)
		for _, frame := range stackArr {
			str += format.TBeg
			str += frame
			str += format.Sep
		}
	}
	return str
}

func (err *ErrRoot) formatJSON(format Format) map[string]interface{} {
	if err == nil {
		return nil
	}
	rootMap := make(map[string]interface{})
	rootMap["message"] = fmt.Sprint(err.Msg)
	if format.WithTrace {
		rootMap["stack"] = formatStackFrames(err.Stack, format.TSep)
	}
	return rootMap
}

// ErrLink represents a single error frame and the accompanying message.
type ErrLink struct {
	Msg   string
	Frame StackFrame
}

func (eLink *ErrLink) formatStr(format Format) string {
	var str string
	str += eLink.Msg
	str += format.Msg
	if format.WithTrace {
		str += format.TBeg
		str += eLink.Frame.format(format.TSep)
	}
	return str
}

func (eLink *ErrLink) formatJSON(format Format) map[string]interface{} {
	wrapMap := make(map[string]interface{})
	wrapMap["message"] = fmt.Sprint(eLink.Msg)
	if format.WithTrace {
		wrapMap["stack"] = eLink.Frame.format(format.TSep)
	}
	return wrapMap
}

func stackContains(stack []StackFrame, frame StackFrame) bool {
	for _, f := range stack {
		if f == frame {
			return true
		}
	}
	return false
}

func formatStackFrames(s []StackFrame, sep string) []string {
	var str []string
	for _, f := range s {
		str = append(str, f.format(sep))
	}
	return str
}
