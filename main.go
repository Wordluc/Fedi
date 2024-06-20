package main

import (
	KeyBoardImpl "TUI/Devices/KeyBoard/impl"
	TermImpl "TUI/Devices/Terminal/impl"
	"TUI/Engine"
	"TUI/Engine/Events"
	"strconv"
)

var keyb = KeyBoardImpl.ImplKeyBoard{}
var t = TermImpl.Terminal{}
var core, e = Engine.Setup(&keyb, &t)

func main() {
	if core == nil {
		panic("Core is nil")
	}
	e = t.Start()
	if e != nil {
		panic(e)
	}
	t.Clear()
	defer t.Stop()
	if e != nil {
		panic(e)
	}
	posEvent:=Events.KeyEvent{
     Key:[]byte{'p'},
     Handler: func(key byte) {
			 x,y:=t.GetCursor()
			t.PrintStr("x: "+strconv.Itoa(x)+" y: "+strconv.Itoa(y))
     },
	}
	core.AddEvent(posEvent)
 	e = core.AddEvent(Events.Arrow(core))
	if e != nil {
		panic(e)
	}
	e = keyb.Start(loop)
	if e != nil {
		panic(e)
	}
	defer keyb.Stop()
}

func loop() bool {
	core.LoopEvent()
	v, _ := keyb.GetKey()
	t.Print([]byte{v})
	if v == 'q' {
		return true
	}
	return false
}
