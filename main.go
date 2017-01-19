package main

import (
	"fmt"

	"github.com/dalloriam/opentis100/tis100"
)

func main() {

	var err error

	program := `@0
	mov up acc
	add acc
	mov acc down
@4
	mov up acc
	mov acc left
	`

	fmt.Println("Input Program: ")

	fmt.Println(program)

	comp := tis100.New("TIS-100")

	err = comp.LoadProgramSource("testProgram", program)

	if err != nil {
		panic(err)
	}

	in := make(chan int)

	go func() {
		for i := 0; i < 50; i++ {
			in <- i
		}
		close(in)
		panic("DOne")
	}()

	err = comp.AttachInput(in, 0, tis100.UP)

	if err != nil {
		panic(err)
	}

	out := make(chan int)

	err = comp.AttachOutput(out, 4, tis100.LEFT)

	if err != nil {
		panic(err)
	}

	comp.Start()

	for o := range out {
		fmt.Print("Output: ")
		fmt.Println(o)
	}
}
