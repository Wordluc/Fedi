package ITerminal

type Size struct {
	Height int
	Width  int
}
type ITerminal interface {
	Start() error
	Stop()
	Clear()
	Print(byte []byte)
	PrintStr(str string)
	Len() Size
}
