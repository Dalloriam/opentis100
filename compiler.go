package opentis100

import (
	"github.com/dalloriam/opentis100/compiler"
	"fmt"
)

type instructionSet struct {
	Labels map[string]int
	Instructions []Instruction
}

type Program struct {
	InstructionSets map[int]*instructionSet
}

func compile(programSource string) (*Program, error) {

	outch, _ := compiler.Tokenize(programSource)

	for tkn := range outch {
		fmt.Println(tkn)
	}

	return nil, nil
}
