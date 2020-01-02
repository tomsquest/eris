package main

// a very basic example under development
import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
)

func compare99(randN int) (err error) {
	if randN > 99 {
		err = eris.New("Integer is greater than 99")
	} else {
		err = nil
	}
	return err
}

func compare100(randN int) <-chan error {
	err := make(chan error, 1)
	go func() {
		if randN > 100 {
			err <- eris.New("Integer is greater than 100")
		} else {
			err <- nil
		}
		close(err)
	}()
	return err
}

func compare99e(randN int) (err error) {
	if randN > 99 {
		err = errors.New("Integer is greater than 99")
	} else {
		err = nil
	}
	return err
}

func compare100e(randN int) <-chan error {
	err := make(chan error, 1)
	go func() {
		if randN > 100 {
			err <- errors.New("Integer is greater than 100")
		} else {
			err <- nil
		}
		close(err)
	}()
	return err
}

func compare101(randN int) <-chan eris.UnpackedError {
	err := make(chan eris.UnpackedError, 1)
	go func() {
		if randN > 101 {
			err <- eris.Unpack(eris.New("Integer is greater than 101"))
		} else {
			err <- eris.UnpackedError{}
		}
		close(err)
	}()
	return err
}

var a = eris.New("global error")

func main() {
	// no wrapping (eris)
	// nErr := compare99(105)
	// fmt.Println(fmt.Sprintf("%+v", nErr))

	// err := compare100(105)
	// fmt.Println(fmt.Sprintf("%+v", <-err))

	// uErr := compare101(105)
	// newErr := <-uErr
	// fmt.Println(newErr.ToString(eris.NewDefaultFormat(true)))

	// no wrapping (pkg/errors)
	// nErr = compare99e(105)
	// fmt.Println(fmt.Sprintf("%+v", nErr))

	// err = compare100e(105)
	// fmt.Println(fmt.Sprintf("%+v", <-err))

	// wrapping (pks/errors)
	err := compare99e(105)
	fmt.Println(fmt.Sprintf("%+v", errors.Wrap(err, "Wrapped pkg/error")))
	fmt.Println()

	// wrapping (pks/errors)
	// err := compare99(105)
	// fmt.Println(fmt.Sprintf("%+v", err))
	// err = eris.Wrap(err, "Wrapped eris root")
	// fmt.Println(fmt.Sprintf("%+v", err))

	// global
	// fmt.Println(fmt.Sprintf("%+v", a))
	// err = eris.Wrap(a, "Wrapped eris global")
	// // err = eris.Wrap(err, "Wrapped eris global again")
	// fmt.Println(fmt.Sprintf("%+v", err))

}
