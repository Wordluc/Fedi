package main

import (
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)

func initCarosello() *Carosello[*TodoBlock, []string] {
	updateCallback := func(display *TodoBlock, data []string) {
		display.SetElement(data[0], data[1])
	}
	newCallback := func(nDisplay int) *TodoBlock {
		return CreateTodoBlock(0, 3*nDisplay, 20)
	}
	selectCallback := func(display *TodoBlock) {
		display.Select()
	}
	deselectCallback := func(display *TodoBlock) {
		display.Deselect()
	}
	callback := Callbacks[*TodoBlock, []string]{
		updateDisplay:   updateCallback,
		newDisplay:      newCallback,
		selectDisplay:   selectCallback,
		deselectDisplay: deselectCallback,
	}

	carosello := CreateCarosello(3, callback)
	return carosello
}

var carosello *Carosello[*TodoBlock, []string]

func main() {
	keyb := Keyboard.Keyboard{}
	term := Terminal.Terminal{}
	core, e := GTUI.NewGtui(loop, &keyb, &term)
	if e != nil {
		panic(e)
	}
	core.SetVisibilityCursor(false)
	wid, hig := core.Size()
	outline := Drawing.CreateRectangleFull(0, 0, wid, hig)
	todos := [][]string{
		{"title1", "text1"},
		{"title2", "text2"},
		{"title3", "text3"},
		{"title4", "text4"},
		{"title5", "text5"},
	}
	carosello = initCarosello()
	carosello.SetPos(2, 1)
	for i := range todos {
		carosello.AddData(todos[i])
	}

	core.AddDrawing(outline)
	core.AddContainer(carosello)
	core.Start()
}

func loop(keyb Keyboard.IKeyBoard, core *GTUI.Gtui) bool {
	if keyb.IsKeySPressed(Keyboard.CtrlQ) {
		return false
	}
	if keyb.IsKeyPressed('l') {
		carosello.Next()
	}
	if keyb.IsKeyPressed('h') {
		carosello.Pre()
	}
	if keyb.IsKeyPressed('d') {
		i, _ := carosello.GetSelectedElement()
		carosello.DeleteData(i)
	}

	return true
}
