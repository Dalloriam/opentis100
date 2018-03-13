package opentis100

import (
	"github.com/dalloriam/opentis100/compiler"
)

func buildProgram(parsed *compiler.Program) *Program {
	outProgram := &Program{
		InstructionSets: map[int]*instructionSet{},
	}

	for i := 0; i < len(parsed.Nodes); i++ {

		set := &instructionSet{
			Labels:       map[string]int{},
			Instructions: make([]*Instruction, len(parsed.Nodes[i].Statements)),
		}

		for i, stmt := range parsed.Nodes[i].Statements {
			if stmt.Label != "" {
				// Label needs to be linked to statement index
				set.Labels[stmt.Label] = i
			}
			set.Instructions[i] = newInstruction(stmt.Instruction.Op, stmt.Instruction.Arguments)
		}

		outProgram.InstructionSets[parsed.Nodes[i].ID] = set
	}

	return outProgram
}

func compile(programSource string) (*Program, error) {

	outch, tknErrors := compiler.Tokenize(programSource)

	p := compiler.NewParser(outch)
	parsedProgram, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}

	for e := range tknErrors {
		return nil, e
	}

	built := buildProgram(parsedProgram)

	return built, nil
}
