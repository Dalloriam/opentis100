package tis100

// SimpleRegister defines a simple register
type SimpleRegister struct {
	value int
}

// NewRegister makes a new simple register
func newRegister() *SimpleRegister {
	r := SimpleRegister{value: 0}

	return &r
}

// Reader returns the reading channel of the register
func (r *SimpleRegister) Read() int {
	return r.value
}

// Writer returns the writing channel of the register
func (r *SimpleRegister) Write(i int) {
	r.value = i
}

type readOnlyRegister struct {
	reader <-chan int
}

func newReadOnlyRegister(inChannel <-chan int) *readOnlyRegister {
	return &readOnlyRegister{reader: inChannel}
}

func (r *readOnlyRegister) Read() int {
	return <-r.reader
}

func (r *readOnlyRegister) Write(i int) {
}

type writeOnlyRegister struct {
	writer chan<- int
}

func newWriteOnlyRegister(outChannel chan<- int) *writeOnlyRegister {
	return &writeOnlyRegister{writer: outChannel}
}

func (r *writeOnlyRegister) Read() int {
	return -1
}

func (r *writeOnlyRegister) Write(i int) {
	r.writer <- i
}

type VirtualRegister struct {
	value int
	ch    chan int
}

func newVirtualRegister() *VirtualRegister {
	r := VirtualRegister{value: 0, ch: make(chan int)}

	return &r
}

func (r *VirtualRegister) Write(i int) {
	r.ch <- i
}

func (r *VirtualRegister) Read() int {
	return <-r.ch
}
