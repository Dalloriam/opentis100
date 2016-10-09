package tis100

type state int

// Different states of the node
const (
	IDLE  state = 0
	RUN   state = 1
	READ  state = 2
	WRITE state = 3
)

type direction int

const (
	up    direction = 0
	right direction = 1
	down  direction = 2
	left  direction = 3
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
	memory             []*instruction

	// Physical Registers
	acc IRegister
	bak IRegister

	// Virtual Registers
	UP    IRegister
	RIGHT IRegister
	DOWN  IRegister
	LEFT  IRegister
}

func newNode(id int) *Node {
	return &Node{ID: id, acc: newRegister(), bak: newRegister(), State: IDLE, ProgramLoaded: false, memory: []*instruction{}}
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
func (n *Node) AttachNode(otherNode *Node, port direction) {
	switch port {
	case up:
		n.UP = otherNode
		otherNode.DOWN = n
	case right:
		n.RIGHT = otherNode
		otherNode.LEFT = n
	case down:
		n.DOWN = otherNode
		otherNode.UP = n
	case left:
		n.LEFT = otherNode
		otherNode.RIGHT = n
	}
}
