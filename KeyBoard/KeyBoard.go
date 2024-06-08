package KeyBoard

import (
	"TUI/Terminal/Token"
)

type IKeyBoard interface {
	Start()error
	Stop()
	GetKey() (byte, error)
	IsPressed(token Token.Token) bool
}

