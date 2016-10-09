package tis100

type state int

// Different states of the node
const (
	IDLE  state = 0
	RUN   state = 1
	READ  state = 2
	WRITE state = 3
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

// CreateNode creates and returns a new node
func CreateNode(id int) *Node {
	return &Node{ID: id, acc: NewRegister(), bak: NewRegister(), State: IDLE, ProgramLoaded: false, memory: []*instruction{}}
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
