package main

import (
	KeyBoardImpl "TUI/Devices/KeyBoard/impl"
	TermImpl "TUI/Devices/Terminal/impl"
	"TUI/Engine"
	"TUI/Engine/Events"
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
