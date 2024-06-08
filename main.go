package main

import (
	"TUI/KeyBoard/impl"
	"TUI/Terminal"
)

var keyb = impl.KeyBoard{}
var t =Terminal.Terminal{}

func main() {
	t.Start()
	t.Clear()
	defer t.Stop()
	e := keyb.Start(loop)
	if e != nil {
		panic(e)
	}
	defer keyb.Stop()
}
func loop() bool {
	v, _ := keyb.GetKey()
	t.Print([]byte{v})
	if v == 'q' {
		return true
	}
	return false
}
