package tis100

import "errors"

const maxNodeX int = 4
const maxNodeY int = 3

// Computer represents a TIS-100 computer
type Computer struct {
	Name  string
	nodes []*Node
}

// New returns a new instance of the TIS-100
func New(computerName string) *Computer {
	nodeList := []*Node{}

	for y := 0; y < maxNodeY; y++ {
		for x := 0; x < maxNodeX; x++ {
			nodeList = append(nodeList, newNode(maxNodeX*y+x))
		}
	}

	for y := 0; y < maxNodeY; y++ {
		for x := 0; x < maxNodeX; x++ {
			node := nodeList[maxNodeX*y+x]

			if 0 < y && y < maxNodeY {
				node.AttachNode(nodeList[maxNodeX*(y-1)+x], UP)
			}

			if 0 < x && x < maxNodeX {
				node.AttachNode(nodeList[maxNodeX*y+x-1], LEFT)
			}
		}
	}

	return &Computer{Name: computerName, nodes: nodeList}
}

// AttachInput attaches an input stream to the specified port of the specified node
func (c *Computer) AttachInput(inStream <-chan int, nodeID int, nodeDirection Direction) error {
	var err error

	r := newReadOnlyRegister(inStream)

	if nodeID < len(c.nodes) {
		err = c.nodes[nodeID].SetPort(nodeDirection, r)
	} else {
		err = errors.New("The specified node doesn't exist!")
	}

	return err
}

// AttachOutput attaches an output stream to the specified port of the specified node
func (c *Computer) AttachOutput(outStream chan<- int, nodeID int, nodeDirection Direction) error {
	var err error

	r := newWriteOnlyRegister(outStream)

	if nodeID < len(c.nodes) {
		err = c.nodes[nodeID].SetPort(nodeDirection, r)
	} else {
		err = errors.New("The specified node doesn't exist!")
	}

	return err
}

// LoadProgramSource compiles the source of a program and loads it into the TIS-100
func (c *Computer) LoadProgramSource(name string, src string) error {
	var err error
	var p *Program

	p, err = Compile(name, src)

	return err
}

// LoadProgramBinary loads a compiled program in the TIS-100
func (c *Computer) LoadProgramBinary(p Program) error {
	var err error

	return err
}
