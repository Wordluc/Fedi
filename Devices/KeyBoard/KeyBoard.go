package KeyBoard

import (
	"TUI/Devices/Terminal/Token"
)

type Loop func() bool
type IKeyBoard interface {
	Start(loop Loop)error
	Stop()
	GetKey() (byte, error)
	IsPressed(token Token.Token) bool
}

