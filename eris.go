package eris

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// New creates a new root error with a static message.
func New(msg string) error {
	return &rootError{
		msg:   msg,
		stack: callers(),
	}
}

// NewGlobal creates a new root error for use as a global sentinel type.
func NewGlobal(msg string) error {
	return &rootError{
		global: true,
		msg:    msg,
		stack:  callers(),
	}
}

// Errorf creates a new root error with a formatted message.
func Errorf(format string, args ...interface{}) error {
	return &rootError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// todo: maybe try to change back to wrap helper and see if the bug is fixable

// Wrap adds additional context to all error types while maintaining the type of the original error.
//
// This method behaves differently for each error type. For root errors, the stack trace is reset to the current
// callers which ensures traces are correct when using global/sentinel error values. Wrapped error types are simply
// wrapped with the new context. For external types (i.e. something other than root or wrap errors), a new root
// error is created for the original error and then it's wrapped with the additional context.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *rootError:
		if e.global {
			e.stack = callers()
		}
	case *wrapError:
	default:
		err = &rootError{
			msg:   e.Error(),
			stack: callers(),
		}
	}

	return &wrapError{
		msg:   msg,
		err:   err,
		frame: caller(),
	}
}

// Wrapf adds additional context to all error types while maintaining the type of the original error.
//
// This is a convenience method for wrapping errors with formatted messages and is otherwise the same as Wrap.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *rootError:
		if e.global {
			e.stack = callers()
		}
	case *wrapError:
	default:
		err = &rootError{
			msg:   e.Error(),
			stack: callers(),
		}
	}

	msg := fmt.Sprintf(format, args...)
	return &wrapError{
		msg:   msg,
		err:   err,
		frame: caller(),
	}
}

// Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method
// returning error. Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if it implements a method
// Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && err == target {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// Cause returns the root cause of the error, which is defined as the first error in the chain. The original
// error is returned if it does not implement `Unwrap() error` and nil is returned if the error is nil.
func Cause(err error) error {
	for {
		uerr := Unwrap(err)
		if uerr == nil {
			return err
		}
		err = uerr
	}
}

type rootError struct {
	global bool
	msg    string
	stack  *stack
}

func (e *rootError) Error() string {
	return fmt.Sprint(e)
}

// todo: document available runes
func (e *rootError) Format(s fmt.State, verb rune) {
	printError(e, s, verb)
}

func (e *rootError) Is(target error) bool {
	if err, ok := target.(*rootError); ok {
		return e.msg == err.msg
	}
	return e.msg == target.Error()
}

type wrapError struct {
	msg   string
	err   error
	frame *frame
}

func (e *wrapError) Error() string {
	return fmt.Sprint(e)
}

func (e *wrapError) Format(s fmt.State, verb rune) {
	printError(e, s, verb)
}

func (e *wrapError) Is(target error) bool {
	if err, ok := target.(*wrapError); ok {
		return e.msg == err.msg
	}
	return e.msg == target.Error()
}

func (e *wrapError) Unwrap() error {
	return e.err
}

func printError(err error, s fmt.State, verb rune) {
	var withTrace, withJSON bool
	if s.Flag('+') {
		withTrace = true
	}
	if s.Flag('#') {
		withJSON = true
	}
	format := NewDefaultFormat(withTrace)
	uErr := Unpack(err)

	var str string
	switch verb {
	case 's':
		fallthrough
	case 'v':
		if withJSON {
			bytes, _ := json.MarshalIndent(uErr.ToJSON(format), "", "\t")
			str = string(bytes)
		} else {
			str = uErr.ToString(format)
		}
	}

	_, _ = io.WriteString(s, str)
}
