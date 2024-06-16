package main

import (
	"TUI/Core"
	KeyImpl "TUI/KeyBoard/impl"
	"TUI/Terminal/Token"
	TermImpl "TUI/Terminal/impl"
)

var keyb = KeyImpl.ImplKeyBoard{}
var t = TermImpl.Terminal{}

func main() {
	t.Start()
	t.Clear()
	defer t.Stop()
	Core.Setup(&keyb, &t)
	eventArrow := Core.ArrowEvent{
		Key: []Token.Token{Token.Arrow_Left, Token.Arrow_Right, Token.Arrow_Up, Token.Arrow_Down},
		Handler: func(Key Token.Token) {
			switch Key {
			case Token.Arrow_Left:
				t.PrintStr("\x1b[1D")
			case Token.Arrow_Right:
				t.PrintStr("\x1b[1C")
			case Token.Arrow_Up:
				t.PrintStr("\x1b[1A")
			case Token.Arrow_Down:
				t.PrintStr("\x1b[1B")
			}
		},
	}
	Core.AddEvent(eventArrow)
	e := keyb.Start(loop)
	if e != nil {
		panic(e)
	}
	defer keyb.Stop()
}
func loop() bool {
	Core.LoopEvent()
	v, _ := keyb.GetKey()
	t.Print([]byte{v})
	if v == 'q' {
		return true
	}
	return false
}
