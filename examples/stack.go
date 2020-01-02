package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
)

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

func main() {
	fmt.Println("\n--- local eris ---")
	localErr := localErr()
	localErr = eris.Wrap(localErr, "new context")
	fmt.Printf("%+v", localErr)
	fmt.Println("------")

	fmt.Println("\n--- nested local eris ---")
	nestedLocalErr := nestedLocalErr()
	nestedLocalErr = eris.Wrap(nestedLocalErr, "new context")
	fmt.Printf("%+v", nestedLocalErr)
	fmt.Println("------")

	fmt.Println("\n--- nested local pkg/errors ---")
	nestedLocalPkgErr := nestedLocalPkgErr()
	nestedLocalPkgErr = errors.Wrap(nestedLocalPkgErr, "new context on pkg error")
	fmt.Printf("%+v", nestedLocalPkgErr)
	fmt.Println("\n------")

	fmt.Println("\n--- global eris ---")
	globalErr := globalErr()
	globalErr = eris.Wrap(globalErr, "new context")
	fmt.Printf("%+v", globalErr)
	fmt.Println("------")

	fmt.Println("\n--- nested global eris ---")
	nestedGlobalErr := nestedGlobalErr()
	nestedGlobalErr = eris.Wrap(nestedGlobalErr, "new context")
	err := eris.Wrap(nestedGlobalErr, "more context")
	fmt.Printf("%+v", err)
	fmt.Println("------")

	fmt.Println("\n--- nested global pkg/errors ---")
	nestedGlobalPkgErr := nestedGlobalPkgErr()
	nestedGlobalPkgErr = errors.Wrap(nestedGlobalPkgErr, "new context")
	err = errors.Wrap(nestedGlobalPkgErr, "more context")
	fmt.Printf("%+v", err)
	fmt.Println("\n------")
}
