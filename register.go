package opentis100

import (
	"errors"
)

type simpleRegister struct {
	value int
}

func newRegister() *simpleRegister {
	r := simpleRegister{value: 0}

	return &r
}

func (r *simpleRegister) Read() (int, error) {
	return r.value, nil
}

func (r *simpleRegister) Write(i int) {
	r.value = i
}

func (r *simpleRegister) Exit() {}

type readOnlyRegister struct {
	reader <-chan int
}

func newReadOnlyRegister(inChannel <-chan int) *readOnlyRegister {
	return &readOnlyRegister{reader: inChannel}
}

func (r *readOnlyRegister) Read() (int, error) {
	val, stillOpen := <-r.reader
	if !stillOpen {
		return -1, errors.New("No more values in register")
	}
	return val, nil
}

func (r *readOnlyRegister) Exit() {}

func (r *readOnlyRegister) Write(i int) {}

type writeOnlyRegister struct {
	writer chan<- int
}

func newWriteOnlyRegister(outChannel chan<- int) *writeOnlyRegister {
	return &writeOnlyRegister{writer: outChannel}
}

func (r *writeOnlyRegister) Read() (int, error) {
	return -1, errors.New("Can't read from a write-only register")
}

func (r *writeOnlyRegister) Exit() {
	close(r.writer)
}

func (r *writeOnlyRegister) Write(i int) {
	r.writer <- i
}

type virtualRegister struct {
	ch     chan int
	isOpen bool
}

func newVirtualRegister() *virtualRegister {
	r := virtualRegister{ch: make(chan int), isOpen: true}

	return &r
}

func (r *virtualRegister) Write(i int) {
	r.ch <- i
}

func (r *virtualRegister) Read() (int, error) {
	val, stillOpen := <-r.ch
	if !stillOpen {
		return -1, errors.New("No more values in register")
	}
	return val, nil
}

func (r *virtualRegister) Exit() {
	if r.isOpen {
		r.isOpen = false
		close(r.ch)
	}
}
