package main

import (
	"fmt"

	"github.com/dalloriam/opentis100/tis100"
)

func main() {

	var err error

	program := `@0
	mov up acc
	mov acc left
	`

	comp := tis100.New("TIS-100")

	fmt.Println("Loading test program in TIS-100...")
	err = comp.LoadProgramSource("testProgram", program)

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting & attaching input goroutine...")
	in := make(chan int)

	go func() {
		for i := 0; i < 50; i++ {
			in <- i
		}
		close(in)
	}()

	err = comp.AttachInput(in, 0, tis100.UP)

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting & attaching output goroutine...")
	out := make(chan int)

	go func() {
		for i := range out {
			fmt.Print("TIS-100: ")
			fmt.Println(i)
		}
	}()

	err = comp.AttachOutput(out, 0, tis100.LEFT)

	if err != nil {
		panic(err)
	}

	comp.Start()
}
