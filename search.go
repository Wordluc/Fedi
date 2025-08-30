package main

import (
	"fmt"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
)

type Search struct {
	*Component.Modal
	textbox      *Component.TextBox
	core         *GTUI.Gtui
	howManyFound *Drawing.TextField
}

func CreateSearch(core *GTUI.Gtui) *Search {
	w, _ := core.Size()
	x := w/2 - 30/2
	modal := Component.CreateModal(30, 5)
	textBox, _ := Component.CreateTextBox(1, 1, 25, 3, core.CreateStreamingCharacter())
	textBox.IsOneLine = true
	howManyFound := Drawing.CreateTextField(26, 2, "")
	title := Drawing.CreateTextField(0, 4, "Search")
	modal.AddComponent(textBox)
	modal.AddDrawing(title, howManyFound)
	modal.SetPos(x, 0)
	modal.SetVisibility(false)
	return &Search{
		Modal:        modal,
		howManyFound: howManyFound,
		textbox:      textBox,
		core:         core,
	}
}

func (s *Search) Close() {
	s.textbox.ClearAll()
	s.textbox.OnLeave()
	s.SetVisibility(false)
	s.SetActive(false)
	s.core.SetVisibilityCursor(false)
	s.howManyFound.SetText("")
}
func (s *Search) Open() {
	s.SetVisibility(true)
	s.SetActive(true)
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

func (s *Search) SetHowManyTODOsFound(number int) {
	s.howManyFound.SetText(fmt.Sprint(number))
}
