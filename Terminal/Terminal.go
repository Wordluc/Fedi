package Terminal

type ITerminal interface {
	Start() error
	Stop()
	Clear()
	Print(byte []byte)
	PrintStr(str string)
}
