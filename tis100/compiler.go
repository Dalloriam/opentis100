package tis100

import (
	"regexp"
	"strconv"
	"strings"
)

// InstructionSet is a set of instructions loaded in a single node
type InstructionSet struct {
	Instructions []instruction

	// Map [labelName]InstructionNo
	Labels map[string]int
}

// ParseBlock parses the instruction set for a single node
func ParseBlock(lines []string) (*InstructionSet, error) {
	var err error
	var ins *InstructionSet

	return ins, err
}

// Program represents a program that can be executed by the TIS-100
type Program struct {
	Name string
	Sets map[int]*InstructionSet // One InstructionSet per node
}

// Compile compile TIS source code and returns an executable program
func Compile(name string, src string) (*Program, error) {
	var err error
	var prog *Program

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
					set, err := ParseBlock(currentBlock)

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
		set, err := ParseBlock(currentBlock)

		if err != nil {
			return nil, err
		}

		sets[currentID] = set
	}

	prog = &Program{Name: name, Sets: sets}
	return prog, err
}
