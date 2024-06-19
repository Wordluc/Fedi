package Engine

import (
	"TUI/Devices/Terminal/Token"
)

type Handler func(Key Token.Token)
type Event interface {
	EventHandler(core *Core) error
}
