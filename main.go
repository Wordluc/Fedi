package main

import (
	"fmt"
	"strings"

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
var viewModal *ViewModal
var searchModal *Search

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
	carosello = initCarosello(wid, hig)
	carosello.SetPos(2, 2)
	edit = CreateEditBlock(wid, hig, 40, 10, core)
	if edit == nil {
		panic("")
	}
	repository = NewRepositoty("prova.csv",
		func(s []string) TODO {
			text := strings.ReplaceAll(s[2], "/n", "\n")
			return TODO{
				Id:     s[0],
				Title:  s[1],
				Text:   text,
				Date:   s[3],
				Status: s[4],
			}
		},
		func(t TODO) []string {
			text := strings.ReplaceAll(t.Text, "\n", "/n")
			return []string{t.Id, t.Title, text, t.Date, t.Status}
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
	viewModal = CreateViewModal(wid, hig, core)
	searchModal = CreateSearch(core)

	core.AddDrawing(outline, title, numberTodos, helper)
	core.AddContainer(carosello)
	core.AddContainer(edit.container)
	core.AddComplexElement(tutorialModal)
	core.AddComplexElement(viewModal)
	core.AddComplexElement(searchModal)
	core.Start()
}

func loop(keyb Keyboard.IKeyBoard, core *GTUI.Gtui) bool {
	if keyb.IsKeySPressed(Keyboard.CtrlQ) {
		editTODO = nil
		if edit.IsOpen() {
			edit.Close()
			return true
		}
		if viewModal.IsOpen() {
			viewModal.Close()
			return true
		}
		if searchModal.IsOpen() {
			searchModal.Close()
			return true
		}
		if tutorialModal.GetVisibility() {
			tutorialModal.SetVisibility(false)
			return true
		}
		return false
	}
	if keyb.IsKeySPressed(Keyboard.Esc) {
		closeAll()
	}
	if edit.IsOpen() {
		if keyb.IsKeySPressed(Keyboard.CtrlK) {
			edit.ActiveTitle()
		} else if keyb.IsKeySPressed(Keyboard.CtrlJ) {
			edit.ActiveText()
		}
		cursorMovement(core, keyb)
	}

	if !edit.IsOpen() {
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
		edit.Open()
		edit.Set(ele.Title, ele.Text)
		editTODO = &ele
	}

	if viewModal.IsOpen() {
		_, ele := carosello.GetSelectedElement()
		viewModal.Change(ele.Title, ele.Text, ele.Status)
	}

	if searchModal.IsOpen() {
		toSearchFor := searchModal.GetText()
		data, e := repository.Get()
		if e != nil {
			return true
		}
		newData := []TODO{}
		for i := range data {
			if ok := strings.Contains(data[i].Title, toSearchFor); ok {
				newData = append(newData, data[i])
				continue
			}
			if ok := strings.Contains(data[i].Text, toSearchFor); ok {
				newData = append(newData, data[i])
				continue
			}
			if ok := strings.Contains(data[i].Status, toSearchFor); ok {
				newData = append(newData, data[i])
				continue
			}
		}
		carosello.Reset()
		carosello.AddDataAll(newData...)
	}

	if keyb.IsKeySPressed(Keyboard.CtrlR) {
		carosello.Reset()
		data, e := repository.Get()
		if e == nil {
			carosello.AddDataAll(data...)
		}
	}

	manageMarksTodos(keyb)
	manageOpenCloseModal(keyb)
	return true
}
