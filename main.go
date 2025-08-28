package main

import (
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)

func initCarosello(width int) *Carosello[*TodoBlock, []string] {
	updateCallback := func(display *TodoBlock, data []string) {
		if len(data) != 2 {
			display.SetElement("", "")
			return
		}
		display.SetElement(data[0], data[1])
	}
	newCallback := func(nDisplay int) *TodoBlock {
		return CreateTodoBlock(0, 3*nDisplay, width-7)
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
var edit *EditBlock

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
	title := Drawing.CreateTextField(0, 0, "TODO")
	title.SetLayer(2)
	todos := [][]string{}
	carosello = initCarosello(wid)
	carosello.SetPos(2, 2)
	edit = CreateEditBlock(wid, hig, 40, 10, core)
	if edit == nil {
		panic("")
	}
	for i := range todos {
		carosello.AddData(todos[i])
	}

	core.AddDrawing(outline, title)
	core.AddContainer(carosello)
	core.AddContainer(edit.container)
	core.Start()
}

var i = 0

func loop(keyb Keyboard.IKeyBoard, core *GTUI.Gtui) bool {
	if keyb.IsKeySPressed(Keyboard.CtrlQ) {
		if edit.IsOn() {
			edit.Toggle()
			core.SetVisibilityCursor(false)
			return true
		}
		return false
	}
	if keyb.IsKeySPressed(Keyboard.CtrlK) {
		if edit.IsOn() {
			edit.ActiveTitle()
		}
	}
	if keyb.IsKeySPressed(Keyboard.CtrlJ) {
		if edit.IsOn() {
			edit.ActiveText()
		}
	}
	if keyb.IsKeyPressed('l') || keyb.IsKeyPressed('j') {
		carosello.Next()
	}
	if keyb.IsKeyPressed('h') || keyb.IsKeyPressed('k') {
		carosello.Pre()
	}
	if keyb.IsKeyPressed('d') {
		i, _ := carosello.GetSelectedElement()
		carosello.DeleteData(i)
	}
	if keyb.IsKeySPressed(Keyboard.CtrlS) {
		isOn := edit.Toggle()
		if !isOn {
			if title, text := edit.GetContent(); text != "" || title != "" {
				GTUI.Log(title)
				carosello.AddData([]string{title, text})
			}
		}
	}

	return true
}
