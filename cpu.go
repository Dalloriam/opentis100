package opentis

import (
	"bytes"
	"encoding/gob"
	"errors"
)

const maxNodeX int = 4
const maxNodeY int = 3

// Computer represents a TIS-100 computer
type Computer struct {
	Name  string
	nodes []*Node
	Debug bool
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
func (c *Computer) CompileProgram(src string) ([]byte, error) {
	var err error
	b, err := compile(src)
	return b, err
}

// LoadProgramBinary loads a compiled program in the TIS-100
func (c *Computer) LoadProgramBinary(b []byte) error {
	var err error

	p := program{}

	buf := bytes.Buffer{}
	buf.Write(b)

	d := gob.NewDecoder(&buf)

	err = d.Decode(&p)

	if err != nil {
		return err
	}

	for i, set := range p.Sets {
		c.nodes[i].LoadInstructions(set)
	}

	return err
}

// Start begins the execution of the currently loaded binary
func (c *Computer) Start() {
	for i := 0; i < len(c.nodes); i++ {
		go c.nodes[i].Start()
	}
}

// Stop stops the execution of the currently loaded binary
func (c *Computer) Stop() {
	for i := 0; i < len(c.nodes); i++ {
		c.nodes[i].Stop()
	}
}