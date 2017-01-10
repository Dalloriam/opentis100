package main

import (
	"fmt"

	"github.com/dalloriam/opentis100/tis100"
)

func main() {
	program := `@0
	mov 0 acc
	mov acc down

	@1
	mov 0 acc
	a: mov 1 acc
	`

	comp := tis100.New("TIS-100")
	comp.LoadProgramSource("testProgram", program)
	fmt.Println(comp.Name)
}
