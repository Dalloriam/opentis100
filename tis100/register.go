package tis100

import "fmt"

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
			r.value = x
			fmt.Printf("Value: %d\n", x)
			if !open {
				close(r.writer)
				return
			}
		}
	}
}

type readOnlyRegister struct {
	reader <-chan int
}

func newReadOnlyRegister(inChannel <-chan int) *readOnlyRegister {
	return &readOnlyRegister{reader: inChannel}
}

func (r *readOnlyRegister) Reader() <-chan int {
	return r.reader
}

func (r *readOnlyRegister) Writer() chan<- int {
	return nil
}

type writeOnlyRegister struct {
	writer chan<- int
}

func newWriteOnlyRegister(outChannel chan<- int) *writeOnlyRegister {
	return &writeOnlyRegister{writer: outChannel}
}

func (r *writeOnlyRegister) Reader() <-chan int {
	return nil
}

func (r *writeOnlyRegister) Writer() chan<- int {
	return r.writer
}
