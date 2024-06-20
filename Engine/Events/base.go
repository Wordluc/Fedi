package Events

import (
	"TUI/Devices/Terminal/Token"
	"TUI/Engine"
)

type TokenEvent struct {
	Key     []Token.Token
	Handler Engine.Handler
}

func (event TokenEvent) EventHandler(core *Engine.Core) error {
	for _, key := range event.Key {
		if core.Keyb.IsTokenPressed(key) {
			event.Handler(byte(key))
		}
	}
	return nil
}
type KeyEvent struct {
	Key     []byte
	Handler Engine.Handler
}

func (event KeyEvent) EventHandler(core *Engine.Core) error {
	for _, key := range event.Key {
		if core.Keyb.IsKeyPressed(key) {
			event.Handler(key)
		}
	}
	return nil
}
