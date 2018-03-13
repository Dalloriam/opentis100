package compiler

type Instruction struct {
	Op        string
	Arguments []string
}

type Statement struct {
	Label       string
	Instruction *Instruction
}

type NodeBlock struct {
	ID         int
	Statements []*Statement
}

type Program struct {
	Nodes map[int]*NodeBlock
}
