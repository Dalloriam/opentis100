package opentis100

type simpleRegister struct {
	value int
}

func newRegister() *simpleRegister {
	r := simpleRegister{value: 0}

	return &r
}

func (r *simpleRegister) Read() int {
	return r.value
}

func (r *simpleRegister) Write(i int) {
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

type virtualRegister struct {
	ch chan int
}

func newVirtualRegister() *virtualRegister {
	r := virtualRegister{ch: make(chan int)}

	return &r
}

func (r *virtualRegister) Write(i int) {
	r.ch <- i
}

func (r *virtualRegister) Read() int {
	return <-r.ch
}
