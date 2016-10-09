package tis100

import "errors"

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
	transfer           chan int
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
	if current == nil {
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

// Reader returns the transfer channel of the node
func (n *Node) Reader() <-chan int {
	return n.transfer
}

// Writer returns the transfer channel of the node
func (n *Node) Writer() chan<- int {
	return n.transfer
}

// LoadBytecode loads bytecode into the node
func (n *Node) LoadBytecode(bytecode []instruction) {
	n.ProgramLoaded = true
}

// UnloadBytecode clears the node's memory and registers
func (n *Node) UnloadBytecode() {
	n.ProgramLoaded = false
	n.memory = []*instruction{}
	n.acc.Writer() <- 0
	n.bak.Writer() <- 0
	n.State = IDLE
}

// AttachNode bidirectionally attaches a node to this node
func (n *Node) AttachNode(otherNode *Node, port Direction) {
	switch port {
	case UP:
		n.up = otherNode
		otherNode.down = n
	case RIGHT:
		n.right = otherNode
		otherNode.left = n
	case DOWN:
		n.down = otherNode
		otherNode.up = n
	case LEFT:
		n.left = otherNode
		otherNode.right = n
	}
}
