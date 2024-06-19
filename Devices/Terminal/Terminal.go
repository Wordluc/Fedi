package ITerminal

type Size struct {
	Width  int
	Height int
}
type ITerminal interface {
	Start() error
	Stop()
	Clear()
	Print(byte []byte)
	PrintStr(str string)
	Len() Size
}
