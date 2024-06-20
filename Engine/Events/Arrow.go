package Events

import (
	"TUI/Devices/Terminal/Token"
	"TUI/Engine"
)

func Arrow(core *Engine.Core) TokenEvent {
	return TokenEvent{
		Key: []Token.Token{Token.Arrow_Left, Token.Arrow_Right, Token.Arrow_Up, Token.Arrow_Down},
		Handler: func(key byte) {
			Key := Token.Token(key)
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
