package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)

func initCarosello(width int) *Carosello[*TodoBlock, TODO] {
	updateCallback := func(display *TodoBlock, data TODO) {
		display.SetElement(data.Title, data.Text, data.Date)
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
	callback := Callbacks[*TodoBlock, TODO]{
		updateDisplay:   updateCallback,
		newDisplay:      newCallback,
		selectDisplay:   selectCallback,
		deselectDisplay: deselectCallback,
	}

	carosello := CreateCarosello(3, callback)
	return carosello
}

var carosello *Carosello[*TodoBlock, TODO]
var edit *EditBlock
var numberTodos *Drawing.TextField
var repository *Repositoty[TODO]

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
	title := Drawing.CreateTextField(0, 0, "TODO:")
	title.SetLayer(2)
	carosello = initCarosello(wid)
	carosello.SetPos(2, 2)
	edit = CreateEditBlock(wid, hig, 40, 10, core)
	if edit == nil {
		panic("")
	}
	repository = NewRepositoty("prova.csv",
		func(s []string) TODO {
			return TODO{
				Id:    s[0],
				Title: s[1],
				Text:  s[2],
				Date:  s[3],
			}
		},
		func(t TODO) []string {
			return []string{t.Id, t.Title, t.Text, t.Date}
		},
		func(t1, t2 TODO) bool {
			return t1.Id == t2.Id
		},
	)
	data, err := repository.Get()
	if err != nil {
		panic(err)
	}
	for i := range data {
		carosello.AddData(data[i])
	}

	numberTodos = Drawing.CreateTextField(5, 0, "0")

	core.AddDrawing(outline, title, numberTodos)
	core.AddContainer(carosello)
	core.AddContainer(edit.container)
	core.Start()
}

var i = 0

func loop(keyb Keyboard.IKeyBoard, core *GTUI.Gtui) bool {
	if keyb.IsKeySPressed(Keyboard.CtrlQ) {
		if edit.IsOn() {
			edit.Toggle()
			return true
		}
		return false
	}
	if keyb.IsKeySPressed(Keyboard.CtrlK) {
		if edit.IsOn() {
			edit.ActiveTitle()
		} else {
			carosello.Pre()
		}
	}
	if keyb.IsKeySPressed(Keyboard.CtrlJ) {
		if edit.IsOn() {
			edit.ActiveText()
		} else {
			carosello.Next()
		}
	}

	if !edit.IsOn() {
		if keyb.IsKeyPressed('l') || keyb.IsKeyPressed('j') {
			carosello.Next()
		}
		if keyb.IsKeyPressed('h') || keyb.IsKeyPressed('k') {
			carosello.Pre()
		}
	}

	if keyb.IsKeySPressed(Keyboard.CtrlD) {
		i, ele := carosello.GetSelectedElement()
		if err := repository.Remove(ele); err == nil {
			carosello.DeleteData(i)
		}
		numberTodos.SetText(fmt.Sprint(len(carosello.GetElements())))
	}
	if keyb.IsKeySPressed(Keyboard.CtrlS) {
		isOn := edit.Toggle()
		if !isOn {
			if title, text := edit.GetContent(); text != "" || title != "" {
				ele := TODO{
					Id:    uuid.NewString(),
					Title: title,
					Text:  text,
					Date:  time.Now().Format(time.RFC850),
				}
				if err := repository.Add(ele); err == nil {
					carosello.AddData(ele)
				}
				numberTodos.SetText(fmt.Sprint(len(carosello.GetElements())))
			}
		}
	}

	return true
}
