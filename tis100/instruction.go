package tis100

type instruction struct {
	Operation string
	Arg1      string
	Arg2      string
}

func newInstruction(op string, arg1 string, arg2 string) *instruction {
	return &instruction{Operation: op, Arg1: arg1, Arg2: arg2}
}
