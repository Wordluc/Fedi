package main

import (
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
)

type EditBlock struct {
	container *Component.Container
	text      *Component.TextBox
	title     *Component.TextBox
	core      *GTUI.Gtui
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
	title, err := Component.CreateTextBox(x+1, y, widBlock-1, 3, core.CreateStreamingCharacter())
	title.IsOneLine = true
	text, err := Component.CreateTextBox(x+1, y+3, widBlock-1, heighetBlock-3, core.CreateStreamingCharacter())
	if err != nil {
		return nil
	}
	container.AddComponent(text)
	container.AddComponent(title)
	container.SetLayer(2)
	container.SetActive(false)
	container.SetVisibility(false)
	return &EditBlock{
		container: container,
		text:      text,
		title:     title,
		core:      core,
	}
}

func (e *EditBlock) Toggle() bool {
	e.container.SetActive(!e.container.GetActivity())
	e.container.SetVisibility(e.container.GetActivity())
	e.core.SetVisibilityCursor(e.container.GetActivity())
	if e.container.GetActivity() {
		e.text.ClearAll()
		e.title.ClearAll()
		e.ActiveTitle()
	} else {
		e.text.OnLeave()
		e.title.OnLeave()
	}
	return e.container.GetActivity()
}
func (e *EditBlock) ActiveText() {
	x, y := e.text.GetPos()
	e.text.OnClick()
	e.title.OnLeave()
	e.core.SetCur(x+1, y+1)
}

func (e *EditBlock) ActiveTitle() {
	x, y := e.title.GetPos()
	e.title.OnClick()
	e.text.OnLeave()
	e.core.SetCur(x+1, y+1)
}

func (e *EditBlock) IsOn() bool {
	return e.container.GetActivity()
}

func (e *EditBlock) GetContent() (string, string) {
	return e.title.GetRawText(), e.text.GetRawText()
}
