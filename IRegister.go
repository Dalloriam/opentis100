package opentis100

type iRegister interface {
	Read() int
	Write(int)
}
