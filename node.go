package opentis100

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	Debug         bool
	State         state
	ProgramLoaded bool

	// Node Execution Information
	currentInstruction int
	running            bool
	memory             *instructionSet

	// Physical Registers
	acc iRegister
	bak iRegister

	// Virtual Registers
	up    iRegister
	right iRegister
	down  iRegister
	left  iRegister
}

func newNode(id int, debug bool) *Node {
	return &Node{ID: id, acc: newRegister(), bak: newRegister(), State: IDLE, ProgramLoaded: false, memory: nil, Debug: debug}
}

func (n *Node) getAllRegisters() []*iRegister {
	return []*iRegister{&n.up, &n.right, &n.down, &n.left, &n.acc, &n.bak}
}

// GetPort returns the register connected at this port
func (n *Node) GetPort(d Direction) *iRegister {
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
func (n *Node) SetPort(d Direction, r iRegister) error {
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
func (n *Node) LoadInstructions(i *instructionSet) {
	n.ProgramLoaded = true
	n.memory = i
}

// UnloadBytecode clears the node's memory and registers
func (n *Node) UnloadBytecode() {
	n.ProgramLoaded = false
	n.memory = &instructionSet{}
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

func (n *Node) getRegister(arg string) (iRegister, error) {

	arg = strings.ToLower(arg)

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

	if arg == "nil" {
		return 0, nil
	}

	reg, err := n.getRegister(arg)

	if err != nil {
		var val int
		val, err = strconv.Atoi(arg)

		if err != nil {
			return 0, err
		}

		return val, nil
	}

	val, err := reg.Read()

	return val, err
}

func (n *Node) tick() error {

	// Do nothing if nothing in memory
	if n.memory == nil {
		return errors.New("Nothing to do")
	}

	// Get current instruction
	ins := n.memory.Instructions[n.currentInstruction]
	n.log(fmt.Sprint(ins))

	// Parse and run instruction
	switch ins.Operation {
	case mov:
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
	case add:
		inValue, err := n.getArgValue(ins.Arg1)

		if err != nil {
			return err
		}

		val, err := n.acc.Read()

		if err != nil {
			return err
		}

		n.acc.Write(val + inValue)
	case sub:
		inValue, err := n.getArgValue(ins.Arg1)

		if err != nil {
			return err
		}

		val, err := n.acc.Read()

		if err != nil {
			return err
		}

		n.acc.Write(val - inValue)
	case nop:
		// Skip instruction (Add 0 to ACC)
		val, err := n.acc.Read()

		if err != nil {
			return err
		}

		n.acc.Write(val + 0)

	case swp:
		// Swap ACC and BAK registers
		tmp, err := n.acc.Read()
		if err != nil {
			return err
		}

		bak, err := n.bak.Read()
		if err != nil {
			return err
		}

		n.acc.Write(bak)
		n.bak.Write(tmp)
	case sav:
		// Copy ACC to BAK register
		val, err := n.acc.Read()
		if err != nil {
			return err
		}
		n.bak.Write(val)

	case neg:
		val, err := n.acc.Read()
		if err != nil {
			return err
		}
		n.acc.Write(-val)

	case jmp:
		tgt := ins.Arg1
		n.currentInstruction = n.memory.Labels[tgt] - 1 // TODO: Errorcheck

	case jez:
		tgt := ins.Arg1

		val, err := n.acc.Read()
		if err != nil {
			return err
		}

		if val == 0 {
			n.currentInstruction = n.memory.Labels[tgt] - 1 // TODO: Errorcheck
		}

	case jlz:
		tgt := ins.Arg1

		val, err := n.acc.Read()
		if err != nil {
			return err
		}

		if val < 0 {
			n.currentInstruction = n.memory.Labels[tgt] - 1 // TODO: Errorcheck
		}

	case jgz:
		tgt := ins.Arg1

		val, err := n.acc.Read()
		if err != nil {
			return err
		}

		if val > 0 {
			n.currentInstruction = n.memory.Labels[tgt] - 1 // TODO: Errorcheck
		}

	default:
		return errors.New("unknown instruction")
	}
	// Update current instruction index
	n.currentInstruction++

	if n.currentInstruction > len(n.memory.Instructions)-1 {
		n.currentInstruction = 0
	}

	return nil
}

func (n *Node) log(msg string) {
	if n.Debug {
		fmt.Printf("[Node %d] - %s\n", n.ID, msg)
	}
}

// Start starts a node
func (n *Node) Start(wg *sync.WaitGroup) {
	var err error

	n.running = true

	for {
		// Stop node if TIS-100 has been stopped
		if !n.running {
			break
		}

		err = n.tick()

		if err != nil {
			n.log(err.Error())
			n.running = false
		}
	}
	n.log("Shutdown")

	for _, reg := range n.getAllRegisters() {
		if *reg != nil {
			(*reg).Exit()
		}
	}
	wg.Done()
}

// Stop stops a node
func (n *Node) Stop() {
	n.running = false
}
