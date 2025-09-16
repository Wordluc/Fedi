package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Keyboard"
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

	var file string
	if len(os.Args) == 1 {
		file = "default"
	} else {
		file = os.Args[1]
	}
	file = "workspace/" + file
	repository = initRepository(file)
	core, e := GTUI.NewGtui(loop)
	if e != nil {
		panic(e)
	}
	defer core.Start()
	core.SetVisibilityCursor(true)
	wid, hig := core.Size()
	outline := Drawing.CreateRectangleFull(0, 0, wid, hig)
	title := Drawing.CreateTextField(0, 0, "TODO:")
	title.SetLayer(2)
	carosello = initCarosello(core, wid, hig)
	carosello.SetPos(2, 2)
	edit = CreateEditBlock(wid, hig, 40, 10, core)
	if edit == nil {
		panic("")
	}
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
	core.AddContainer(carosello, edit)
	core.AddComplexElement(tutorialModal, viewModal, searchModal)
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
			core.GoToComponent(edit.titleField)
		} else if keyb.IsKeySPressed(Keyboard.CtrlJ) {
			edit.ActiveText()
			core.GoToComponent(edit.textField)
		}
		cursorMovement(core, keyb)
	}
	if searchModal.IsOpen() {
		cursorMovement(core, keyb)
	}
	if !edit.IsOpen() && !searchModal.IsOpen() {
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
	if keyb.IsKeySPressed(Keyboard.CtrlV) {
		ele := core.GetCurrentComponent()
		if textBox, ok := ele.(*Component.TextBox); ok {
			textBox.Paste(keyb.GetClickboard())
		}
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
			if ok := strings.Contains(strings.ToLower(data[i].Title), strings.ToLower(toSearchFor)); ok {
				newData = append(newData, data[i])
				continue
			}
			if ok := strings.Contains(strings.ToLower(data[i].Text), strings.ToLower(toSearchFor)); ok {
				newData = append(newData, data[i])
				continue
			}
			if ok := strings.Contains(strings.ToLower(data[i].Status), strings.ToLower(toSearchFor)); ok {
				newData = append(newData, data[i])
				continue
			}
			if ok := strings.Contains(strings.ToLower(data[i].Date), strings.ToLower(toSearchFor)); ok {
				newData = append(newData, data[i])
				continue
			}
		}
		carosello.Reset()
		carosello.AddDataAll(newData...)
		searchModal.SetHowManyTODOsFound(len(newData))
	}

	if keyb.IsKeySPressed(Keyboard.CtrlR) {
		carosello.Reset()
		data, e := repository.Get()
		if e == nil {
			carosello.AddDataAll(data...)
		}
	}

	manageMarksTodos(keyb)
	manageOpenCloseModal(core, keyb)
	return true
}
