package Engine

import (
)

type Handler func(Key byte)
type Event interface {
	EventHandler(core *Core) error
}
