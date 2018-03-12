package opentis100

import (
	"errors"
	"sync"
)

const maxNodeX int = 4
const maxNodeY int = 3

// Computer represents a TIS-100 computer
type Computer struct {
	Name  string
	nodes []*Node
	Debug bool
	wg    *sync.WaitGroup
}

// New returns a new instance of the TIS-100
func New(computerName string, debug bool) *Computer {
	nodeList := []*Node{}

	for y := 0; y < maxNodeY; y++ {
		for x := 0; x < maxNodeX; x++ {
			nodeList = append(nodeList, newNode(maxNodeX*y+x, debug))
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

	return &Computer{Name: computerName, nodes: nodeList, Debug: debug}
}

// AttachInput attaches an input stream to the specified port of the specified node
func (c *Computer) AttachInput(inStream <-chan int, nodeID int, nodeDirection Direction) error {
	var err error

	r := newReadOnlyRegister(inStream)

	if nodeID < len(c.nodes) {
		err = c.nodes[nodeID].SetPort(nodeDirection, r)
	} else {
		err = errors.New("the specified node doesn't exist")
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
		err = errors.New("the specified node doesn't exist")
	}

	return err
}

// CompileProgram compiles the source of a program and loads it into the TIS-100
func (c *Computer) CompileProgram(src string) (*Program, error) {
	ins, err := compile(src)
	return ins, err
}

// LoadProgramBinary loads a compiled program in the TIS-100
func (c *Computer) LoadProgram(p *Program) error {
	var err error

	for i, set := range p.InstructionSets {
		c.nodes[i].LoadInstructions(set)
	}

	return err
}

// Start begins the execution of the currently loaded binary
func (c *Computer) Start() {
	var wg sync.WaitGroup

	for i := 0; i < len(c.nodes); i++ {
		go c.nodes[i].Start(&wg)
		wg.Add(1)
	}
	c.wg = &wg
}

func (c *Computer) WaitUntilComplete() {
	c.wg.Wait()
}

// Stop stops the execution of the currently loaded binary
func (c *Computer) Stop() {
	for i := 0; i < len(c.nodes); i++ {
		c.nodes[i].Stop()
	}
}
