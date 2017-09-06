package opentis

type iRegister interface {
	Read() (int, error)
	Write(int)
	Exit()
}
