package main

import (
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
)

type Search struct {
	*Component.Modal
	textbox *Component.TextBox
	core    *GTUI.Gtui
}

func CreateSearch(core *GTUI.Gtui) *Search {
	w, _ := core.Size()
	x := w/2 - 30/2
	modal := Component.CreateModal(30, 5)
	textBox, _ := Component.CreateTextBox(1, 1, 28, 3, core.CreateStreamingCharacter())
	textBox.IsOneLine = true
	title := Drawing.CreateTextField(0, 4, "Search")
	modal.AddComponent(textBox)
	modal.AddDrawing(title)
	modal.SetPos(x, 0)
	modal.SetVisibility(false)
	return &Search{
		Modal:   modal,
		textbox: textBox,
		core:    core,
	}
}

func (s *Search) Close() {
	s.textbox.ClearAll()
	s.textbox.OnLeave()
	s.SetVisibility(false)
	s.core.SetVisibilityCursor(false)
}
func (s *Search) Open() {
	s.SetVisibility(true)
	s.textbox.OnClick()
	x, y := s.textbox.GetPos()
	s.core.SetVisibilityCursor(true)
	s.core.SetCur(x+1, y+1)
}

func (s *Search) GetText() string {
	return s.textbox.GetText()
}

func (s *Search) IsOpen() bool {
	return s.GetVisibility()
}
