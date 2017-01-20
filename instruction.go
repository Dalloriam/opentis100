package opentis100

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
)

// Instruction represents a single TIS-100 instruction
type Instruction struct {
	Operation operation
	Arg1      string
	Arg2      string
}

func newInstruction(op string, arg1 string, arg2 string) *Instruction {
	var o operation

	switch op {
	case "mov":
		o = mov
	case "add":
		o = add
	case "sub":
		o = sub
	case "nop":
		o = nop
	case "SAV":
		o = sav
	case "NEG":
		o = neg
	}
	return &Instruction{Operation: o, Arg1: arg1, Arg2: arg2}
}
