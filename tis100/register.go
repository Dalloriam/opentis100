package tis100

// SimpleRegister defines a simple register
type SimpleRegister struct {
	value  int
	reader chan int
	writer chan int
}

// NewRegister makes a new simple register
func newRegister() *SimpleRegister {
	r := SimpleRegister{value: 0, reader: make(chan int), writer: make(chan int)}

	go r.valueUpdater()
	go r.valueReader()

	return &r
}

// Reader returns the reading channel of the register
func (r *SimpleRegister) Reader() <-chan int {
	return r.reader
}

// Writer returns the writing channel of the register
func (r *SimpleRegister) Writer() chan<- int {
	return r.writer
}

func (r *SimpleRegister) valueReader() {
	for {
		r.reader <- r.value
	}
}

func (r *SimpleRegister) valueUpdater() {
	for {
		select {
		case x, open := <-r.writer:
			if !open {
				close(r.writer)
				return
			}
			r.value = x
		}
	}
}
