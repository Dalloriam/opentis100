package opentis100

import (
	"strings"
)

type operation int

// All instruction actions
const (
	mov operation = 0

	add operation = 1
	sub operation = 2

	nop operation = 3

	swp operation = 4
	sav operation = 5

	neg operation = 6

	jez operation = 7
	jmp operation = 8
	jlz operation = 9
	jgz operation = 10
)

// Instruction represents a single TIS-100 instruction
type Instruction struct {
	Operation operation
	Arg1      string
	Arg2      string
}

func newInstruction(op string, arg1 string, arg2 string) *Instruction {
	var o operation

	op = strings.ToLower(op)

	switch op {
	case "mov":
		o = mov
	case "add":
		o = add
	case "sub":
		o = sub
	case "nop":
		o = nop
	case "sav":
		o = sav
	case "swp":
		o = swp
	case "neg":
		o = neg
	case "jez":
		o = jez
	case "jmp":
		o = jmp
	case "jlz":
		o = jlz
	case "jgz":
		o = jgz
	}
	return &Instruction{Operation: o, Arg1: arg1, Arg2: arg2}
}
