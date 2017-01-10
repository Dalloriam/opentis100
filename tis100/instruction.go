package tis100

// Instruction represents a single TIS-100 instruction
type Instruction struct {
	Operation string
	Arg1      string
	Arg2      string
}

func newInstruction(op string, arg1 string, arg2 string) *Instruction {
	return &Instruction{Operation: op, Arg1: arg1, Arg2: arg2}
}
