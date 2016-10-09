package main

import (
	"fmt"

	"github.com/dalloriam/opentis/tis100"
)

func main() {
	program := ` @0
	mov 1 acc
	`

	comp := tis100.New("TIS-100")
	comp.LoadProgram(program)
	fmt.Println(comp.Name)
}
