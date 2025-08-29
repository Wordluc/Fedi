package main

import (
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
)

type EditBlock struct {
	container  *Component.Container
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
	x := widScreen/2 - widBlock/2
	y := heightScreen - heighetBlock
	container := Component.CreateContainer(x, y)
	outline := Drawing.CreateRectangle(x, y, widBlock, heighetBlock-1)
	titleField, err := Component.CreateTextBox(x+1, y+1, widBlock-2, 3, core.CreateStreamingCharacter())
	titleField.IsOneLine = true
	textField, err := Component.CreateTextBox(x+1, y+4, widBlock-2, heighetBlock-6, core.CreateStreamingCharacter())
	if err != nil {
		return nil
	}
	titleModal := Drawing.CreateTextField(x, y, "Edit")
	container.AddDrawing(titleModal)
	container.AddDrawing(outline)
	container.AddComponent(textField)
	container.AddComponent(titleField)
	container.SetLayer(2)
	container.SetActive(false)
	container.SetVisibility(false)
	return &EditBlock{
		container:  container,
		textField:  textField,
		titleModal: titleModal,
		titleField: titleField,
		core:       core,
	}
}

func (e *EditBlock) Toggle(isOn bool) bool {
	e.container.SetActive(isOn)
	e.container.SetVisibility(isOn)
	e.core.SetVisibilityCursor(isOn)
	if e.container.GetActivity() {
		e.textField.ClearAll()
		e.titleField.ClearAll()
		e.ActiveTitle()
	} else {
		e.textField.OnLeave()
		e.titleField.OnLeave()
	}
	return e.container.GetActivity()
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

func (e *EditBlock) IsOn() bool {
	return e.container.GetActivity()
}

func (e *EditBlock) GetContent() (string, string) {
	return e.titleField.GetText(), e.textField.GetText()
}

func (e *EditBlock) Set(title, text string) {
	e.textField.ClearAll()
	e.textField.Paste(text)
	e.titleField.ClearAll()
	e.titleField.Paste(title)
}

func (e *EditBlock) SetTitleModal(text string) {
	e.titleModal.SetText(text)
}
