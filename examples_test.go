package eris_test

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
)

// todo: change all line numbers to what they should be in the go playground?
//       or maybe reorganize these examples into whole file examples so the line numbers will be correct automatically?
//       this would be unpack_test.go, local_test.go, global_test.go, and external_test.go?

// todo: add Tests to run examples that don't have predictable output
//       it'll be useful to see this output in CI to verify that submitted PRs haven't messed up the output

// Demonstrates unpacking an external (non-eris) error and printing the raw output.
func ExampleUnpack_external() {
	// example func that returns an IO error
	readFile := func(fname string) error {
		return io.ErrUnexpectedEOF
	}

	// unpack and print the raw error
	err := readFile("example.json")
	uerr := eris.Unpack(err)
	fmt.Println(uerr)
	// Output:
	// {{ []} [] unexpected EOF}
}

// Demonstrates unpacking a wrapped error and printing the raw output.
func ExampleUnpack_wrapped() {
	// example func that returns an eris error
	readFile := func(fname string) error {
		return eris.New("unexpected EOF")
	}

	// example func that catches an error and wraps it with additional context
	parseFile := func(fname string) error {
		// read the file
		err := readFile(fname)
		if err != nil {
			return eris.Wrapf(err, "error reading file '%v'", fname)
		}
		return nil
	}

	// unpack and print the raw error
	err := parseFile("example.json")
	uerr := eris.Unpack(err)
	jsonErr, _ := json.MarshalIndent(uerr, "", "\t")
	fmt.Println(string(jsonErr))
}

// Hack to run examples that don't have a predictable output (i.e. all
// examples that involve printing stack traces).
func TestExampleUnpack_wrapped(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleUnpack_wrapped()
}

// Demonstrates JSON formatting of wrapped errors that originate from
// external (non-eris) error types.
func ExampleUnpackedError_ToJSON_external() {
	// example func that returns an IO error
	readFile := func(fname string) error {
		return io.ErrUnexpectedEOF
	}

	// unpack and print the error
	err := readFile("example.json")
	uerr := eris.Unpack(err)
	format := eris.NewDefaultFormat(false) // false: omit stack trace
	u, _ := json.Marshal(uerr.ToJSON(format))
	fmt.Println(string(u))
	// Output:
	// {"external":"unexpected EOF"}
}

// Demonstrates JSON formatting of wrapped errors that originate from
// global root errors (created via eris.NewGlobal).
func ExampleUnpackedError_ToJSON_global() {
	// declare a "global" error type
	ErrUnexpectedEOF := eris.NewGlobal("unexpected EOF")

	// example func that wraps a global error value
	readFile := func(fname string) error {
		return eris.Wrapf(ErrUnexpectedEOF, "error reading file '%v'", fname)
	}

	// example func that catches and returns an error without modification
	parseFile := func(fname string) error {
		// read the file
		err := readFile(fname)
		if err != nil {
			return err
		}
		return nil
	}

	// call parseFile and catch the error
	err := parseFile("example.json")

	// print the error via fmt.Printf
	fmt.Printf("%v\n", err) // %v: omit stack trace

	// example output:
	// unexpected EOF: error reading file 'example.json'

	// unpack and print the error via uerr.ToString(...)
	uerr := eris.Unpack(err)
	format := eris.NewDefaultFormat(true) // true: include stack trace
	fmt.Printf("%v\n", uerr.ToString(format))

	// example output:
	// unexpected EOF
	//   eris_test.ExampleUnpackedError_ToString_global.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 122
	//   eris_test.ExampleUnpackedError_ToString_global.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 128
	//   eris_test.ExampleUnpackedError_ToString_global: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 136
	//   eris_test.TestExampleUnpackedError_ToString_global: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 154
	// error reading file 'example.json'
	//   eris_test.ExampleUnpackedError_ToString_global.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 122
}

// Hack to run examples that don't have a predictable output (i.e. all
// examples that involve printing stack traces).
func TestExampleUnpackedError_ToJSON_global(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleUnpackedError_ToJSON_global()
}

// Demonstrates JSON formatting of wrapped errors that originate from
// local root errors (created at the source of the error via eris.New).
func ExampleUnpackedError_ToJSON_local() {
	fmt.Println("testing")
}

func TestExampleUnpackedError_ToJSON_local(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleUnpackedError_ToJSON_local()
}

// Demonstrates string formatting of wrapped errors that originate from
// external (non-eris) error types.
func ExampleUnpackedError_ToString_external() {
	// example func that returns an IO error
	readFile := func(fname string) error {
		return io.ErrUnexpectedEOF
	}

	// unpack and print the error
	err := readFile("example.json")
	uerr := eris.Unpack(err)
	format := eris.NewDefaultFormat(false) // false: omit stack trace
	fmt.Println(uerr.ToString(format))
	// Output:
	// unexpected EOF
}

// Demonstrates string formatting of wrapped errors that originate from
// global root errors (created via eris.NewGlobal).
func ExampleUnpackedError_ToString_global() {
	// declare a "global" error type
	ErrUnexpectedEOF := eris.NewGlobal("unexpected EOF")

	// example func that wraps a global error value
	readFile := func(fname string) error {
		return eris.Wrapf(ErrUnexpectedEOF, "error reading file '%v'", fname)
	}

	// example func that catches and returns an error without modification
	parseFile := func(fname string) error {
		// read the file
		err := readFile(fname)
		if err != nil {
			return err
		}
		return nil
	}

	// call parseFile and catch the error
	err := parseFile("example.json")

	// print the error via fmt.Printf
	fmt.Printf("%v\n", err) // %v: omit stack trace

	// example output:
	// unexpected EOF: error reading file 'example.json'

	// unpack and print the error via uerr.ToString(...)
	uerr := eris.Unpack(err)
	format := eris.NewDefaultFormat(true) // true: include stack trace
	fmt.Printf("%v\n", uerr.ToString(format))

	// example output:
	// unexpected EOF
	//   eris_test.ExampleUnpackedError_ToString_global.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 122
	//   eris_test.ExampleUnpackedError_ToString_global.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 128
	//   eris_test.ExampleUnpackedError_ToString_global: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 136
	//   eris_test.TestExampleUnpackedError_ToString_global: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 154
	// error reading file 'example.json'
	//   eris_test.ExampleUnpackedError_ToString_global.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 122
}

func TestExampleUnpackedError_ToString_global(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleUnpackedError_ToString_global()
}

// Demonstrates string formatting of wrapped errors that originate from
// local root errors (created at the source of the error via eris.New).
func ExampleUnpackedError_ToString_local() {
	// example func that returns an eris error
	readFile := func(fname string) error {
		return eris.New("unexpected EOF")
	}

	// example func that catches an error and wraps it with additional context
	parseFile := func(fname string) error {
		// read the file
		err := readFile(fname)
		if err != nil {
			return eris.Wrapf(err, "error reading file '%v'", fname)
		}
		return nil
	}

	// unpack and print the error
	err := parseFile("example.json")
	uerr := eris.Unpack(err)
	format := eris.NewDefaultFormat(true) // true: include stack trace
	fmt.Println(uerr.ToString(format))

	// example output:
	//
}

func TestExampleUnpackedError_ToString_local(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleUnpackedError_ToString_local()
}

var GlobalPkgErr = errors.New("new global pkg error")

var GlobalErr = eris.NewGlobal("new global eris error")

func localErr() error {
	return eris.New("new local eris error")
}

func nestedLocalErr() error {
	return localErr()
}

func localPkgErr() error {
	return errors.New("new local pkg error")
}

func nestedLocalPkgErr() error {
	return localPkgErr()
}

func globalErr() error {
	return eris.Wrap(GlobalErr, "some context")
}

func nestedGlobalErr() error {
	return globalErr()
}

func globalPkgErr() error {
	return errors.Wrap(GlobalPkgErr, "some context")
}

func nestedGlobalPkgErr() error {
	return globalPkgErr()
}

// todo: after implementing stack trace unpacking, try hard to break it
//       pretty sure it's possible if there are stack frames between the end of the stack and newer wrap frames

func main() {
	localErr := localErr()
	localErr = eris.Wrap(localErr, "new context")
	fmt.Println("--- local eris string ('+v') ---")
	fmt.Printf("%+v\n\n", localErr)
	fmt.Println("--- local eris JSON ('+#v') ---")
	fmt.Printf("%#+v\n\n", localErr)

	nestedLocalErr := nestedLocalErr()
	nestedLocalErr = eris.Wrap(nestedLocalErr, "new context")
	fmt.Println("--- nested local eris string ('+v') ---")
	fmt.Printf("%+v\n\n", nestedLocalErr)
	fmt.Println("--- nested local eris JSON ('+#v') ---")
	fmt.Printf("%+#v\n\n", nestedLocalErr)

	nestedLocalPkgErr := nestedLocalPkgErr()
	nestedLocalPkgErr = errors.Wrap(nestedLocalPkgErr, "new context on pkg error")
	fmt.Println("--- nested local pkg/errors string ('+v') ---")
	fmt.Printf("%+v\n\n", nestedLocalPkgErr)

	globalErr := globalErr()
	globalErr = eris.Wrap(globalErr, "new context")
	fmt.Println("--- global eris string ('+v') ---")
	fmt.Printf("%+v\n\n", globalErr)
	fmt.Println("--- global eris JSON ('+#v') ---")
	fmt.Printf("%+#v\n\n", globalErr)

	nestedGlobalErr := nestedGlobalErr()
	nestedGlobalErr = eris.Wrap(nestedGlobalErr, "new context")
	err := eris.Wrap(nestedGlobalErr, "more context")
	fmt.Println("--- nested global eris string ('+v') ---")
	fmt.Printf("%+v\n\n", nestedGlobalErr)
	fmt.Println("--- nested global eris JSON ('#v') ---")
	fmt.Printf("%+v\n\n", nestedGlobalErr)
	fmt.Println("--- nested global eris JSON ('+#v') ---")
	fmt.Printf("%+#v\n\n", nestedGlobalErr)

	nestedGlobalPkgErr := nestedGlobalPkgErr()
	nestedGlobalPkgErr = errors.Wrap(nestedGlobalPkgErr, "new context")
	err = errors.Wrap(nestedGlobalPkgErr, "more context")
	fmt.Println("--- nested global pkg/errors ---")
	fmt.Printf("%+v\n\n", err)
}
