package opentis100

import (
	"bytes"
	"encoding/gob"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// InstructionSet is a set of instructions loaded in a single node
type InstructionSet struct {
	Instructions []*Instruction

	// Map [labelName]InstructionNo
	Labels map[string]int
}

func parseInstruction(instruction string) (*Instruction, error) {
	var err error

	data := strings.Split(instruction, " ")

	if len(data) > 0 && len(data) < 4 {
		// Valid instruction
		op := data[0]

		var a1 string
		var a2 string

		if len(data) > 1 {
			a1 = data[1]
		}

		if len(data) > 2 {
			a2 = data[2]
		}

		return newInstruction(op, a1, a2), err

	}
	// Invalid instruction
	return nil, errors.New("[Error] Invalid instruction: " + instruction)

}

func parseBlock(lines []string) (*InstructionSet, error) {
	var err error
	ins := &InstructionSet{Labels: make(map[string]int), Instructions: []*Instruction{}}

	for i, line := range lines {

		// Label detection
		// S: MOV LEFT ACC
		// JEZ E
		// SWP
		// ADD 1
		labels := strings.Split(line, ":")

		if len(labels) > 1 {
			// Line has label
			ins.Labels[labels[0]] = i
		}

		instr := strings.TrimSpace(strings.Replace(line, labels[0]+":", "", -1))

		in, err := parseInstruction(instr)

		if err != nil {
			return nil, err
		}

		ins.Instructions = append(ins.Instructions, in)

	}

	return ins, err
}

type program struct {
	Name string
	Sets map[int]*InstructionSet // One InstructionSet per node
}

func compile(src string) ([]byte, error) {
	var err error
	var prog program

	sets := make(map[int]*InstructionSet)

	commentPattern := regexp.MustCompile("#.*$")

	lines := strings.Split(src, "\n")

	var currentID int
	var currentBlock []string

	for _, line := range lines {

		line = strings.TrimSpace(commentPattern.ReplaceAllString(line, ""))

		if len(line) > 0 {
			if strings.HasPrefix(line, "@") {
				if currentBlock != nil {
					set, err := parseBlock(currentBlock)

					if err != nil {
						return nil, err
					}

					sets[currentID] = set
					currentBlock = nil
				}
				currentID, err = strconv.Atoi(strings.Replace(line, "@", "", -1))
				if err != nil {
					return nil, err
				}
			} else {
				currentBlock = append(currentBlock, line)
			}
		}
	}

	// Push last block
	if currentBlock != nil {
		set, err := parseBlock(currentBlock)

		if err != nil {
			return nil, err
		}

		sets[currentID] = set
	}

	prog = program{Sets: sets}

	// Dump compiled program to binary
	b := bytes.Buffer{}
	encoder := gob.NewEncoder(&b)

	err = encoder.Encode(prog)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), err
}
