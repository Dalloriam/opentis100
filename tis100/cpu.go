package tis100

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
				node.AttachNode(nodeList[maxNodeX*(y-1)+x], up)
			}

			if 0 < x && x < maxNodeX {
				node.AttachNode(nodeList[maxNodeX*y+x-1], left)
			}
		}
	}

	return &Computer{Name: computerName, nodes: nodeList}
}
