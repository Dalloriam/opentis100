package opentis100

type instructionSet struct {
	Labels       map[string]int
	Instructions []*Instruction
}

type Program struct {
	InstructionSets map[int]*instructionSet
}
