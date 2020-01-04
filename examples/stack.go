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
	fmt.Println("--- local eris string ('%+v') ---")
	fmt.Printf("%+v\n\n", localErr)
	fmt.Println("--- local eris JSON ('%#+v') ---")
	fmt.Printf("%#+v\n\n", localErr)

	nestedLocalErr := nestedLocalErr()
	nestedLocalErr = eris.Wrap(nestedLocalErr, "new context")
	fmt.Println("--- nested local eris string ('%+v') ---")
	fmt.Printf("%+v\n\n", nestedLocalErr)
	fmt.Println("--- nested local eris JSON ('%#+v') ---")
	fmt.Printf("%+#v\n\n", nestedLocalErr)

	nestedLocalPkgErr := nestedLocalPkgErr()
	nestedLocalPkgErr = errors.Wrap(nestedLocalPkgErr, "new context on pkg error")
	fmt.Println("--- nested local pkg/errors string ('%+v') ---")
	fmt.Printf("%+v\n\n", nestedLocalPkgErr)

	globalErr := globalErr()
	globalErr = eris.Wrap(globalErr, "new context")
	fmt.Println("--- global eris string ('%+v') ---")
	fmt.Printf("%+v\n\n", globalErr)
	fmt.Println("--- global eris JSON ('%#+v') ---")
	fmt.Printf("%+#v\n\n", globalErr)

	nestedGlobalErr := nestedGlobalErr()
	nestedGlobalErr = eris.Wrap(nestedGlobalErr, "new context")
	err := eris.Wrap(nestedGlobalErr, "more context")
	fmt.Println("--- nested global eris string ('%+v') ---")
	fmt.Printf("%+v\n\n", nestedGlobalErr)
	fmt.Println("--- nested global eris JSON ('%#v') ---")
	fmt.Printf("%+v\n\n", nestedGlobalErr)
	fmt.Println("--- nested global eris JSON ('%#+v') ---")
	fmt.Printf("%+#v\n\n", nestedGlobalErr)

	nestedGlobalPkgErr := nestedGlobalPkgErr()
	nestedGlobalPkgErr = errors.Wrap(nestedGlobalPkgErr, "new context")
	err = errors.Wrap(nestedGlobalPkgErr, "more context")
	fmt.Println("--- nested global pkg/errors ---")
	fmt.Printf("%+v\n\n", err)
}
