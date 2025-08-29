package main

import (
	"fmt"
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)

var carosello *Carosello[*TodoBlock, TODO]
var edit *EditBlock
var numberTodos *Drawing.TextField
var repository *Repositoty[TODO]
var editTODO *TODO
var tutorialModal *Component.Modal

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
				Id:     s[0],
				Title:  s[1],
				Text:   s[2],
				Date:   s[3],
				Status: s[4],
			}
		},
		func(t TODO) []string {
			return []string{t.Id, t.Title, t.Text, t.Date, t.Status}
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
	tutorialModal = CreateTutorialModal(wid, hig)
	tutorialModal.SetVisibility(false)
	tutorialModal.SetActive(false)
	numberTodos = Drawing.CreateTextField(5, 0, "0")
	numberTodos.SetText(fmt.Sprint(len(carosello.GetElements())))
	helper := Drawing.CreateTextField(1, hig-2, "Tab: to open/close tutorial")

	core.AddDrawing(outline, title, numberTodos, helper)
	core.AddContainer(carosello)
	core.AddContainer(edit.container)
	core.AddComplexElement(tutorialModal)
	core.Start()
}

func loop(keyb Keyboard.IKeyBoard, core *GTUI.Gtui) bool {
	if keyb.IsKeySPressed(Keyboard.CtrlQ) {
		if edit.IsOn() {
			edit.Toggle(false)
			editTODO = nil
			return true
		}
		return false
	}
	if edit.IsOn() {
		if keyb.IsKeySPressed(Keyboard.Esc) {
			edit.Toggle(false)
		} else if keyb.IsKeySPressed(Keyboard.CtrlK) {
			edit.ActiveTitle()
		} else if keyb.IsKeySPressed(Keyboard.CtrlJ) {
			edit.ActiveText()
		}
		cursorMovement(core, keyb)
	}

	if !edit.IsOn() {
		if keyb.IsKeyPressed('l') || keyb.IsKeyPressed('j') || keyb.IsKeySPressed(Keyboard.Down) {
			carosello.Next()
		} else if keyb.IsKeyPressed('h') || keyb.IsKeyPressed('k') || keyb.IsKeySPressed(Keyboard.Up) {
			carosello.Pre()
		}
	}
	if keyb.IsKeySPressed(Keyboard.Tab) {
		tutorialModal.SetActive(!tutorialModal.GetActivity())
		tutorialModal.SetVisibility(tutorialModal.GetActivity())
	}

	if keyb.IsKeySPressed(Keyboard.CtrlE) {
		_, ele := carosello.GetSelectedElement()
		edit.SetTitleModal("Edit")
		edit.Toggle(true)
		edit.Set(ele.Title, ele.Text)
		editTODO = &ele
	}

	if keyb.IsKeySPressed(Keyboard.CtrlS) {
		isOn := edit.Toggle(!edit.IsOn())
		if !isOn {
			saveContentEditBlock()
		} else {
			edit.SetTitleModal("New Todo")
		}
	}
	manageMarksTodos(keyb)

	return true
}
