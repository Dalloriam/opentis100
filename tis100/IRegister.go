package tis100

// IRegister describes a virtual register
type IRegister interface {
	Reader() <-chan int
	Writer() chan<- int
}
