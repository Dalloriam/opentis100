package opentis100

// IRegister describes a virtual register
type IRegister interface {
	Read() int
	Write(int)
}
