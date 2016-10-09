package tis100

// Computer represents a TIS-100 computer
type Computer struct {
	Name  string
	nodes []string
}

// New returns a new instance of the TIS-100
func New(inputChannel <-chan int, outputChannel <-chan int) *Computer {
	return &Computer{}
}
