package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/google/uuid"
)

func initCarosello(width, height int) *Carosello[*TodoBlock, TODO] {
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

	carosello := CreateCarosello((height/4)-1, callback)
	return carosello
}
func initRepository(fileName string) *Repositoty[TODO] {
	p, err := os.Executable()
	if err != nil {
		panic(err)
	}
	p = filepath.Dir(p)
	fullpath := path.Join(p, fileName+".csv")
	os.MkdirAll(filepath.Dir(fullpath), os.ModePerm)
	repository = NewRepositoty(fullpath,
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
	return repository
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
	var ele TODO
	if keyb.IsKeySPressed(Keyboard.CtrlD) {
		_, ele = carosello.GetSelectedElement()
		ele.Status = Done
	}
	if keyb.IsKeySPressed(Keyboard.CtrlX) {
		_, ele = carosello.GetSelectedElement()
		ele.Status = Deleted
	}

	if keyb.IsKeySPressed(Keyboard.CtrlA) {
		_, ele = carosello.GetSelectedElement()
		ele.Status = Archived
	}

	if keyb.IsKeySPressed(Keyboard.CtrlW) {
		_, ele = carosello.GetSelectedElement()
		ele.Status = WaitingFor
	}

	if keyb.IsKeySPressed(Keyboard.CtrlP) {
		_, ele = carosello.GetSelectedElement()
		ele.Status = Progress
	}

	if ele.Status != "" && repository.Set(ele) == nil {
		updateCaroselloData()
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
	title, text := edit.GetContent()
	if text == "" && title == "" {
		return
	}
	if editTODO != nil {
		editTODO.Title = title
		editTODO.Text = text
		editTODO.Status = Ready
		if repository.Set(*editTODO) == nil {
			updateCaroselloData()
		}
		editTODO = nil
	} else {
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

func setIconMark(mark *Drawing.TextBlock, icon string) {
	mark.SetText("")
	switch icon {
	case Ready:
		mark.SetText("⛟")
		mark.SetColor(Color.Get(Color.Gray, Color.None))
	case Done:
		mark.SetText("✓")
		mark.SetColor(Color.Get(Color.Green, Color.None))
	case Progress:
		mark.SetText("⛏")
		mark.SetColor(Color.Get(Color.Cyan, Color.None))
	case Deleted:
		mark.SetText("⛔")
	case Archived:
		mark.SetText("X")
		mark.SetColor(Color.Get(Color.Yellow, Color.None))
	case WaitingFor:
		mark.SetText("⚠")
		mark.SetColor(Color.Get(Color.Yellow, Color.None))
	}
}

func closeAll() {
	viewModal.Close()
	searchModal.Close()
	edit.Close()
	tutorialModal.SetVisibility(false)
}
func manageOpenCloseModal(keyb Keyboard.IKeyBoard) {
	if keyb.IsKeySPressed(Keyboard.CtrlV) {
		_, ele := carosello.GetSelectedElement()
		if viewModal.IsOpen() {
			viewModal.Close()
		} else {
			closeAll()
			viewModal.Open(ele.Title, ele.Text, ele.Status)
		}
	}

	if keyb.IsKeySPressed(Keyboard.CtrlF) {
		if searchModal.IsOpen() {
			searchModal.Close()
		} else {
			closeAll()
			searchModal.Open()
		}
	}

	if keyb.IsKeySPressed(Keyboard.CtrlS) {
		isOn := edit.IsOpen()
		if !isOn {
			closeAll()
			edit.SetTitleModal("New Todo")
			edit.Open()
		} else {
			edit.Close()
			saveContentEditBlock()
		}
	}
	if keyb.IsKeySPressed(Keyboard.CtrlE) {
		closeAll()
		_, ele := carosello.GetSelectedElement()
		edit.SetTitleModal("Edit")
		edit.Open()
		edit.Set(ele.Title, ele.Text)
		editTODO = &ele
	}
}
