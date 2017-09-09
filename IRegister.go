package opentis100

type iRegister interface {
	Read() (int, error)
	Write(int)
	Exit()
}
