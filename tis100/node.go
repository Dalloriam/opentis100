package tis100

import (
	"errors"
	"fmt"
	"strconv"
)

type state int

// Different states of the node
const (
	IDLE  state = 0
	RUN   state = 1
	READ  state = 2
	WRITE state = 3
)

// Direction is an alias for one of the node's ports
type Direction int

// All 4 directions
const (
	UP    Direction = 0
	RIGHT Direction = 1
	DOWN  Direction = 2
	LEFT  Direction = 3
)

// Node represents a TIS-100 CPU core
type Node struct {
	// Node Information
	ID            int
	State         state
	ProgramLoaded bool

	// Node Execution Information
	currentInstruction int
	running            bool
	memory             *InstructionSet

	// Physical Registers
	acc IRegister
	bak IRegister

	// Virtual Registers
	up    IRegister
	right IRegister
	down  IRegister
	left  IRegister
}

func newNode(id int) *Node {
	return &Node{ID: id, acc: newRegister(), bak: newRegister(), State: IDLE, ProgramLoaded: false, memory: nil}
}

// GetPort returns the register connected at this port
func (n *Node) GetPort(d Direction) *IRegister {
	switch d {
	case UP:
		return &n.up

	case RIGHT:
		return &n.right

	case DOWN:
		return &n.down

	case LEFT:
		return &n.left
	}
	return nil
}

// SetPort binds a register to a port of this node
func (n *Node) SetPort(d Direction, r IRegister) error {
	current := n.GetPort(d)
	if *current == nil {
		switch d {
		case UP:
			n.up = r

		case RIGHT:
			n.right = r

		case DOWN:
			n.down = r

		case LEFT:
			n.left = r
		}
		return nil
	}
	return errors.New("The specified port is not available")
}

// LoadInstructions loads bytecode into the node
func (n *Node) LoadInstructions(i *InstructionSet) {
	n.ProgramLoaded = true
	n.memory = i
}

// UnloadBytecode clears the node's memory and registers
func (n *Node) UnloadBytecode() {
	n.ProgramLoaded = false
	n.memory = &InstructionSet{}
	n.acc.Write(0)
	n.bak.Write(0)
	n.State = IDLE
}

// AttachNode bidirectionally attaches a node to this node
func (n *Node) AttachNode(otherNode *Node, port Direction) {

	buffer := newVirtualRegister()

	switch port {
	case UP:
		n.up = buffer
		otherNode.down = buffer
	case RIGHT:
		n.right = buffer
		otherNode.left = buffer
	case DOWN:
		n.down = buffer
		otherNode.up = buffer
	case LEFT:
		n.left = buffer
		otherNode.right = buffer
	}
}

func (n *Node) getRegister(arg string) (IRegister, error) {

	switch arg {
	case "up":
		return n.up, nil
	case "down":
		return n.down, nil
	case "right":
		return n.right, nil
	case "left":
		return n.left, nil
	case "acc":
		return n.acc, nil
	default:
		return nil, fmt.Errorf("Register %s does not exist", arg)
	}
}

func (n *Node) getArgValue(arg string) (int, error) {
	reg, err := n.getRegister(arg)

	if err != nil {
		val, err := strconv.Atoi(arg)

		if err != nil {
			return 0, err
		}

		return val, nil
	}

	return reg.Read(), err
}

func (n *Node) tick() error {

	// Do nothing if nothing in memory
	if n.memory == nil {
		return nil
	}

	// Get current instruction
	ins := n.memory.Instructions[n.currentInstruction]

	// Parse and run instruction
	switch ins.Operation {
	case MOV:
		// Read from ins.Arg1, write to ins.Arg2
		inValue, err := n.getArgValue(ins.Arg1)

		if err != nil {
			return err
		}

		outReg, err := n.getRegister(ins.Arg2)
		if err != nil {
			return err
		}

		outReg.Write(inValue)
	case ADD:
		inValue, err := n.getArgValue(ins.Arg1)

		if err != nil {
			return err
		}

		n.acc.Write(n.acc.Read() + inValue)
	case SUB:
		inValue, err := n.getArgValue(ins.Arg1)

		if err != nil {
			return err
		}

		n.acc.Write(n.acc.Read() - inValue)
	default:
		return errors.New("Unknown instruction.")
	}
	// Update current instruction index
	n.currentInstruction++

	if n.currentInstruction > len(n.memory.Instructions)-1 {
		n.currentInstruction = 0
	}

	return nil
}

// Start starts a node
func (n *Node) Start() {
	var err error

	n.running = true

	for {
		// Stop node if TIS-100 has been stopped
		if !n.running {
			break
		}

		err = n.tick()

		if err != nil {
			fmt.Printf("[Node %d] - %s", n.ID, err)
		}
	}
	fmt.Printf("[Node %d] - Shutdown\n", n.ID)
}

// Stop stops a node
func (n *Node) Stop() {
	n.running = false
}
