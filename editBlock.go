package main

import (
	"strings"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
)

type EditBlock struct {
	*Component.Container
	textField  *Component.TextBox
	titleField *Component.TextBox
	titleModal *Drawing.TextField
	core       *GTUI.Gtui
}

func CreateEditBlock(widScreen, heightScreen, widBlock, heighetBlock int, core *GTUI.Gtui) *EditBlock {
	if widBlock > widScreen {
		return nil
	}
	if heighetBlock > heightScreen {
		return nil
	}
	container := Component.CreateContainer()
	outline := Drawing.CreateRectangle(0, 0, widBlock, heighetBlock-1)
	titleField, err := Component.CreateTextBox(1, 1, widBlock-2, 3, core.CreateStreamingCharacter())
	titleField.IsOneLine = true
	textField, err := Component.CreateTextBox(1, 4, widBlock-2, heighetBlock-6, core.CreateStreamingCharacter())
	if err != nil {
		return nil
	}
	titleModal := Drawing.CreateTextField(0, 0, "Edit")
	container.AddDrawing(titleModal, outline)
	container.AddComponent(textField, titleField)
	container.SetActive(false)
	container.SetVisibility(false)
	container.SetLayer(2)

	x := widScreen/2 - widBlock/2
	y := heightScreen - heighetBlock
	container.SetPos(x, y)
	return &EditBlock{
		Container:  container,
		textField:  textField,
		titleModal: titleModal,
		titleField: titleField,
		core:       core,
	}
}
func (e *EditBlock) Close() {
	e.SetActive(false)
	e.SetVisibility(false)
	e.core.SetVisibilityCursor(false)
	e.textField.OnLeave()
	e.titleField.OnLeave()
}

func (e *EditBlock) Open() {
	e.SetActive(true)
	e.SetVisibility(true)
	e.core.SetVisibilityCursor(true)
	e.textField.ClearAll()
	e.titleField.ClearAll()
	e.ActiveTitle()
}
func (e *EditBlock) IsOn() bool {
	return e.GetVisibility()
}
func (e *EditBlock) ActiveText() {
	x, y := e.textField.GetPos()
	e.textField.OnClick()
	e.titleField.OnLeave()
	e.core.SetCur(x+1, y+1)
}

func (e *EditBlock) ActiveTitle() {
	x, y := e.titleField.GetPos()
	e.titleField.OnClick()
	e.textField.OnLeave()
	e.core.SetCur(x+1, y+1)
}

func (e *EditBlock) IsOpen() bool {
	return e.GetActivity()
}

func (e *EditBlock) GetContent() (string, string) {
	GTUI.Log(e.textField.GetText())
	return e.titleField.GetText(), e.textField.GetText()
}

func (e *EditBlock) Set(title, text string) {
	e.textField.ClearAll()
	e.textField.Paste(strings.ReplaceAll(text, "/n", "\n"))
	e.titleField.ClearAll()
	e.titleField.Paste(title)
}

func (e *EditBlock) SetTitleModal(text string) {
	e.titleModal.SetText(text)
}
