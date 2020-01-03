package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
)

// todo: add ExampleLocal, ExampleGlobal, ExamplePkgErrors
//       use closures to demonstrate error wrapping, etc
//       make each example self-contained for godocs
//       integrate play for golang playground

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
	fmt.Println("\n--- local eris string ('%+v') ---")
	fmt.Println(fmt.Sprintf("%+v", localErr))
	fmt.Println("------")
	fmt.Println("\n--- local eris JSON ('%#+j') ---")
	fmt.Println(fmt.Sprintf("%#+j", localErr))
	fmt.Println("------")

	nestedLocalErr := nestedLocalErr()
	nestedLocalErr = eris.Wrap(nestedLocalErr, "new context")
	fmt.Println("\n--- nested local eris string ('%+v') ---")
	fmt.Fprintln("%+v", nestedLocalErr)
	fmt.Println("------")
	fmt.Println("\n--- nested local eris JSON ('%#+j') ---")
	fmt.Fprintln("%+#v", nestedLocalErr)
	fmt.Println("------")

	fmt.Println("\n--- nested local pkg/errors ---")
	nestedLocalPkgErr := nestedLocalPkgErr()
	nestedLocalPkgErr = errors.Wrap(nestedLocalPkgErr, "new context on pkg error")
	fmt.Printf("%+v", nestedLocalPkgErr)
	fmt.Println("\n------")

	globalErr := globalErr()
	globalErr = eris.Wrap(globalErr, "new context")
	fmt.Println("--- global eris string ('%+v') ---")
	fmt.Fprintln("%+v\n", globalErr)
	fmt.Println("--- global eris JSON ('%#+j') ---")
	fmt.Fprintln("%+#j\n\n", globalErr)

	fmt.Println("\n--- nested global eris ---")
	nestedGlobalErr := nestedGlobalErr()
	nestedGlobalErr = eris.Wrap(nestedGlobalErr, "new context")
	err := eris.Wrap(nestedGlobalErr, "more context")
	fmt.Println("\n--- nested global eris string ('%+v') ---")
	fmt.Println(fmt.Sprintf("%+v", nestedGlobalErr))
	fmt.Println("------")
	fmt.Println("\n--- nested global eris JSON ('%+j') ---")
	fmt.Println(fmt.Sprintf("%+j", nestedGlobalErr))
	fmt.Println("------")
	fmt.Println("\n--- nested global eris JSON ('%#+j') ---")
	fmt.Println(fmt.Sprintf("%+#j", nestedGlobalErr))
	fmt.Println("------")

	fmt.Println("\n--- nested global pkg/errors ---")
	nestedGlobalPkgErr := nestedGlobalPkgErr()
	nestedGlobalPkgErr = errors.Wrap(nestedGlobalPkgErr, "new context")
	err = errors.Wrap(nestedGlobalPkgErr, "more context")
	fmt.Printf("%+v", err)
	fmt.Println("\n------")
}
