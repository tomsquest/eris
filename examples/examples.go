package main

// a very basic example under development
import (
	"fmt"

	"github.com/rotisserie/eris"
)

func compare100(randN int) <-chan error {
	err := make(chan error, 1)
	if randN > 100 {
		err <- eris.New("Integer is greater than 100")
	}
	return err
}

func main() {
	err := compare100(105)
	fmt.Println(fmt.Sprintf("%+v", <-err))
}
