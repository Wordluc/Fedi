package Events

import (
	"TUI/Devices/Terminal/Token"
	"TUI/Engine"
)

type ArrowEvent struct {
	Key     []Token.Token
	Handler Engine.Handler
}

func (event ArrowEvent) EventHandler(core *Engine.Core) error {
	for _, key := range event.Key {
		if core.Keyb.IsPressed(key) {
			event.Handler(key)
		}
	}
	return nil
}
//func moveLeft(core *Engine.Core) bool {
//	size:=core.Term.Len()
//	
//
//}
func Arrow(core *Engine.Core) ArrowEvent {
	return ArrowEvent{
		Key: []Token.Token{Token.Arrow_Left, Token.Arrow_Right, Token.Arrow_Up, Token.Arrow_Down},
		Handler: func(Key Token.Token) {
			switch Key {
			case Token.Arrow_Left:
				core.Term.PrintStr("\x1b[1D")
			case Token.Arrow_Right:
				core.Term.PrintStr("\x1b[1C")
			case Token.Arrow_Up:
				core.Term.PrintStr("\x1b[1A")
			case Token.Arrow_Down:
				core.Term.PrintStr("\x1b[1B")
			}
		},
	}

}
