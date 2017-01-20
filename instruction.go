package opentis100

type operation int

// All instruction actions
const (
	MOV operation = 0

	ADD operation = 1
	SUB operation = 2

	NOP operation = 3

	SWP operation = 4
	SAV operation = 5

	NEG operation = 6
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
		o = MOV
	case "add":
		o = ADD
	case "sub":
		o = SUB
	case "nop":
		o = NOP
	case "SAV":
		o = SAV
	case "NEG":
		o = NEG
	}
	return &Instruction{Operation: o, Arg1: arg1, Arg2: arg2}
}
