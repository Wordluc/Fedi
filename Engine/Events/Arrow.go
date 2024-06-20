package Events

import (
	"TUI/Devices/Terminal/Token"
	"TUI/Engine"
)

func Arrow(core *Engine.Core) TokenEvent {
	return TokenEvent{
		Key: []Token.Token{Token.Arrow_Left, Token.Arrow_Right, Token.Arrow_Up, Token.Arrow_Down},
		Handler: func(key byte) {
			xP,yP := core.Term.GetCursor()
      xSize,ySize := core.Term.Len()
			Key := Token.Token(key)
			switch Key {
			case Token.Arrow_Left:
				if xP == 0 {
					return ;
				}
				core.Term.PrintStr("\x1b[1D")
			case Token.Arrow_Right:
				if xP == xSize - 1 {
					return ;
				}
				core.Term.PrintStr("\x1b[1C")
			case Token.Arrow_Up:
				if yP == 1 {
					return ;
				}
				core.Term.PrintStr("\x1b[1A")
			case Token.Arrow_Down:
				if yP == ySize - 1 {
					return ;
				}
				core.Term.PrintStr("\x1b[1B")
			}
		},
	}
}
