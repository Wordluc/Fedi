package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)

func initCarosello(width int) *Carosello[*TodoBlock, TODO] {
	updateCallback := func(display *TodoBlock, data TODO) {
		display.SetElement(data.Title, data.Text, data.Date, data.Status, data.Id)
	}
	newCallback := func(nDisplay int) *TodoBlock {
		return CreateTodoBlock(0, 4*nDisplay, width-7)
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

	carosello := CreateCarosello(10, callback)
	return carosello
}

var carosello *Carosello[*TodoBlock, TODO]
var edit *EditBlock
var numberTodos *Drawing.TextField
var repository *Repositoty[TODO]
var editingTODO *TODO
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
			editingTODO = nil
			return true
		}
		return false
	}
	if edit.IsOn() {
		if keyb.IsKeySPressed(Keyboard.CtrlK) {
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
		editingTODO = &ele
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

func updateCaroselloData() {
	data, err := repository.Get()
	if err != nil {
		return
	}
	carosello.Refresh(data...)
	numberTodos.SetText(fmt.Sprint(len(carosello.GetElements())))
}

func manageMarksTodos(keyb Keyboard.IKeyBoard) {
	if keyb.IsKeySPressed(Keyboard.CtrlD) {
		_, ele := carosello.GetSelectedElement()
		ele.Status = Done
		if repository.Set(ele) == nil {
			updateCaroselloData()
		}
	}
	if keyb.IsKeySPressed(Keyboard.CtrlX) {
		_, ele := carosello.GetSelectedElement()
		ele.Status = Deleted
		if repository.Set(ele) == nil {
			updateCaroselloData()
		}
	}

	if keyb.IsKeySPressed(Keyboard.CtrlA) {
		_, ele := carosello.GetSelectedElement()
		ele.Status = Archived
		if repository.Set(ele) == nil {
			updateCaroselloData()
		}
	}

	if keyb.IsKeySPressed(Keyboard.CtrlW) {
		_, ele := carosello.GetSelectedElement()
		ele.Status = WaitingFor
		if repository.Set(ele) == nil {
			updateCaroselloData()
		}
	}
}

func cursorMovement(core *GTUI.Gtui, keyb Keyboard.IKeyBoard) {
	var x, y = core.GetCur()
	if keyb.IsKeySPressed(Keyboard.Down) {
		y++
	}
	if keyb.IsKeySPressed(Keyboard.Up) {
		y--
	}
	if keyb.IsKeySPressed(Keyboard.Right) {
		x++
	}
	if keyb.IsKeySPressed(Keyboard.Left) {
		x--
	}
	core.SetCur(x, y)
}

func saveContentEditBlock() {
	if editingTODO != nil {
		if title, text := edit.GetContent(); text != "" || title != "" {
			editingTODO.Title = title
			editingTODO.Text = text
			editingTODO.Status = Ready
			if repository.Set(*editingTODO) == nil {
				updateCaroselloData()
			}
			editingTODO = nil
		}
	} else if title, text := edit.GetContent(); text != "" || title != "" {
		text = strings.ReplaceAll(text, "\n", ";")
		ele := TODO{
			Id:     uuid.NewString(),
			Title:  title,
			Text:   text,
			Date:   time.Now().Format("Mon, 02 Jan 2006"),
			Status: Ready,
		}
		if err := repository.Add(ele); err == nil {
			carosello.AddData(ele)
		}
		numberTodos.SetText(fmt.Sprint(len(carosello.GetElements())))
	}
}
