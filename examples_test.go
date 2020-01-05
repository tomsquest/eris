package eris_test

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/rotisserie/eris"
)

// todo: add expected output (anonymized)
// todo: change all line numbers to what they should be in the go playground?

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

	// unpack and print the error via uerr.ToJSON(...)
	err := parseFile("example.json")
	uerr := eris.Unpack(err)
	format := eris.NewDefaultFormat(true) // true: include stack trace
	u, _ := json.MarshalIndent(uerr.ToJSON(format), "", "\t")
	fmt.Printf("%v\n", string(u))

	// example output:
	// {
	//   "root": {
	//     "message": "unexpected EOF",
	//     "stack": [
	//       "eris_test.ExampleUnpackedError_ToJSON_global.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 41",
	//       "eris_test.ExampleUnpackedError_ToJSON_global.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 47",
	//       "eris_test.ExampleUnpackedError_ToJSON_global: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 55",
	//       "eris_test.TestExampleUnpackedError_ToJSON_global: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 87"
	//     ]
	//   },
	//   "wrap": [
	//     {
	//       "message": "error reading file 'example.json'",
	//       "stack": "eris_test.ExampleUnpackedError_ToJSON_global.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 41"
	//     }
	//   ]
	// }
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

	// example func that just catches and returns an error
	processFile := func(fname string) error {
		// parse the file
		err := parseFile(fname)
		if err != nil {
			return err
		}
		return nil
	}

	// another example func that catches and wraps an error
	printFile := func(fname string) error {
		// process the file
		err := processFile(fname)
		if err != nil {
			return eris.Wrapf(err, "error printing file '%v'", fname)
		}
		return nil
	}

	// unpack and print the raw error
	err := printFile("example.json")
	uerr := eris.Unpack(err)
	format := eris.NewDefaultFormat(true) // true: include stack trace
	u, _ := json.MarshalIndent(uerr.ToJSON(format), "", "\t")
	fmt.Printf("%v\n", string(u))

	// example output:
	// {
	//   "root": {
	//     "message": "unexpected EOF",
	//     "stack": [
	//       "eris_test.ExampleUnpackedError_ToJSON_local.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 95",
	//       "eris_test.ExampleUnpackedError_ToJSON_local.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 101",
	//       "eris_test.ExampleUnpackedError_ToJSON_local.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 103",
	//       "eris_test.ExampleUnpackedError_ToJSON_local.func3: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 111",
	//       "eris_test.ExampleUnpackedError_ToJSON_local.func4: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 121",
	//       "eris_test.ExampleUnpackedError_ToJSON_local.func4: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 123",
	//       "eris_test.ExampleUnpackedError_ToJSON_local: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 129",
	//       "eris_test.TestExampleUnpackedError_ToJSON_local: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 167"
	//     ]
	//   },
	//   "wrap": [
	//     {
	//       "message": "error reading file 'example.json'",
	//       "stack": "eris_test.ExampleUnpackedError_ToJSON_local.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 103"
	//     },
	//     {
	//       "message": "error printing file 'example.json'",
	//       "stack": "eris_test.ExampleUnpackedError_ToJSON_local.func4: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 123"
	//     }
	//   ]
	// }
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
	//   eris_test.ExampleUnpackedError_ToString_global.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 195
	//   eris_test.ExampleUnpackedError_ToString_global.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 201
	//   eris_test.ExampleUnpackedError_ToString_global: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 209
	//   eris_test.TestExampleUnpackedError_ToString_global: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 236
	// error reading file 'example.json'
	//   eris_test.ExampleUnpackedError_ToString_global.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 195
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

	// call parseFile and catch the error
	err := parseFile("example.json")

	// print the error via fmt.Printf
	fmt.Printf("%v\n", err) // %v: omit stack trace

	// example output:
	// unexpected EOF: error reading file 'example.json'

	// unpack and print the error via uerr.ToString(...)
	uerr := eris.Unpack(err)
	format := eris.NewDefaultFormat(true) // true: include stack trace
	fmt.Println(uerr.ToString(format))

	// example output:
	// unexpected EOF
	//   eris_test.ExampleUnpackedError_ToString_local.func1: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 244
	//   eris_test.ExampleUnpackedError_ToString_local.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 250
	//   eris_test.ExampleUnpackedError_ToString_local.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 252
	//   eris_test.ExampleUnpackedError_ToString_local: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 258
	//   eris_test.TestExampleUnpackedError_ToString_local: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 286
	// error reading file 'example.json'
	//   eris_test.ExampleUnpackedError_ToString_local.func2: /Users/morningvera/go/src/github.com/rotisserie/eris/examples_test.go: 252
}

func TestExampleUnpackedError_ToString_local(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleUnpackedError_ToString_local()
}
