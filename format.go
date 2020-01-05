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
func (e *UnpackedError) ToString(format Format) string {
	// todo: clean up these conditionals if possible
	var str string
	if len(e.ErrRoot.Msg) != 0 || len(e.ErrRoot.Stack) != 0 {
		str += e.ErrRoot.formatStr(format)
	}
	if format.WithTrace && len(e.ErrChain) != 0 {
		str += format.Sep
	}
	for _, eLink := range e.ErrChain {
		if !format.WithTrace {
			str += format.Sep
		}
		str += eLink.formatStr(format)
	}
	if e.ExternalErr != "" {
		str += fmt.Sprint(e.ExternalErr)
	}
	return str
}

// ToJSON returns a JSON formatted map for a given eris error.
func (e *UnpackedError) ToJSON(format Format) map[string]interface{} {
	// todo: clean up these conditionals if possible
	jsonMap := make(map[string]interface{})
	if len(e.ErrRoot.Msg) != 0 || len(e.ErrRoot.Stack) != 0 {
		jsonMap["root"] = e.ErrRoot.formatJSON(format)
	}

	if len(e.ErrChain) != 0 {
		var wrapArr []map[string]interface{}
		for _, eLink := range e.ErrChain {
			wrapMap := eLink.formatJSON(format)
			wrapArr = append(wrapArr, wrapMap)
		}
		jsonMap["wrap"] = wrapArr
	}

	if e.ExternalErr != "" {
		jsonMap["external"] = fmt.Sprint(e.ExternalErr)
	}

	return jsonMap
}

func (e *UnpackedError) unpackRootErr(err *rootError) {
	e.ErrRoot.Msg = err.msg
	e.ErrRoot.Stack = err.stack.get()
}

func (e *UnpackedError) unpackWrapErr(err *wrapError) {
	// prepend links to match the stack trace order
	link := ErrLink{Msg: err.msg}
	stack := err.stack.get()
	if len(stack) > 0 {
		link.Frame = stack[0]
	}
	e.ErrChain = append([]ErrLink{link}, e.ErrChain...)

	nextErr := err.Unwrap()
	switch nextErr.(type) {
	case *rootError:
		e.unpackRootErr(nextErr.(*rootError))
	case *wrapError:
		e.unpackWrapErr(nextErr.(*wrapError))
	default:
		return
	}

	// combine the wrap stack with the root stack
	e.ErrRoot.Stack.combineStack(stack)
}

// ErrRoot represents an error stack and the accompanying message.
type ErrRoot struct {
	Msg   string
	Stack Stack
}

func (err *ErrRoot) formatStr(format Format) string {
	str := err.Msg
	str += format.Msg
	if format.WithTrace {
		stackArr := err.Stack.format(format.TSep)
		for i, frame := range stackArr {
			str += format.TBeg
			str += frame
			if i < len(stackArr)-1 {
				str += format.Sep
			}
		}
	}
	return str
}

func (err *ErrRoot) formatJSON(format Format) map[string]interface{} {
	rootMap := make(map[string]interface{})
	rootMap["message"] = fmt.Sprint(err.Msg)
	if format.WithTrace {
		rootMap["stack"] = err.Stack.format(format.TSep)
	}
	return rootMap
}

// ErrLink represents a single error frame and the accompanying message.
type ErrLink struct {
	Msg   string
	Frame StackFrame
}

func (eLink *ErrLink) formatStr(format Format) string {
	str := eLink.Msg
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
