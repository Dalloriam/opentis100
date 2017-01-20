package main

import (
	"fmt"

	"github.com/dalloriam/opentis100/tis100"
)

func main() {

	var err error

	program := `@0
	mov up right
@1
	mov left acc
	add right
	mov acc down
@2
	mov up left
@5
	mov up left

@4
	mov right left
	`

	fmt.Println("Input Program: ")

	fmt.Println(program)

	comp := tis100.New("TIS-100")

	err = comp.LoadProgramSource("testProgram", program)

	if err != nil {
		panic(err)
	}

	in1 := make(chan int)
	in2 := make(chan int)

	go func() {
		for i := 0; i < 50; i += 2 {
			in1 <- i
		}
		close(in1)
	}()

	go func() {
		for i := 1; i < 50; i += 2 {
			in2 <- i
		}
		close(in2)
	}()

	err = comp.AttachInput(in1, 0, tis100.UP)
	err = comp.AttachInput(in2, 2, tis100.UP)

	if err != nil {
		panic(err)
	}

	out := make(chan int)

	err = comp.AttachOutput(out, 4, tis100.LEFT)

	if err != nil {
		panic(err)
	}

	comp.Start()

	i := 0
	for o := range out {
		fmt.Print("Output: ")
		fmt.Println(o)
		i++
		if i >= 50 {
			comp.Stop()
			break
		}
	}
}
