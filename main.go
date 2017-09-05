package main

import (
	"fmt"

	"github.com/dalloriam/opentis100/opentis"
)

const program = `@0
mov up acc
mov acc down

@4
mov up acc
add acc
mov acc left
`

func main() {

	inputs := []int{3, 4, 5}
	inChan := make(chan int)

	go func() {
		for _, inpt := range inputs {
			inChan <- inpt
		}
		close(inChan)
	}()

	outChan := make(chan int)

	fmt.Println("Setting up machine...")
	tis := opentis.New("TIS-100", true)
	tis.AttachInput(inChan, 0, opentis.UP)
	tis.AttachOutput(outChan, 4, opentis.LEFT)

	fmt.Println("Compiling test program...")
	compiled, err := tis.CompileProgram(program)

	if err != nil {
		panic(err)
	}

	fmt.Println("Loading test program in memory...")
	err = tis.LoadProgramBinary(compiled)

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting execution...")
	tis.Start()

	for o := range outChan {
		fmt.Printf("Output: %d\n", o)
	}
}
