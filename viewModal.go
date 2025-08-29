package main

import (
	"fmt"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
)

type ViewModal struct {
	*Component.Modal
	textbox *Drawing.TextBlock
	title   *Drawing.TextField
	status  *Drawing.TextField
}

func CreateViewModal(w, h int, core *GTUI.Gtui) *ViewModal {
	const sizeX, sizeY = 90, 30
	modal := Component.CreateModal(sizeX, sizeY)
	x := w/2 - sizeX/2
	y := h/2 - sizeY/2
	title := Drawing.CreateTextField(sizeX/2-4, 1, "View ToDo")
	status := Drawing.CreateTextField(sizeX/2+5, 1, "")
	titleField := Drawing.CreateTextField(2, 2, "")
	title.SetColor(Color.Get(Color.Red, Color.None))
	textbox := Drawing.CreateTextBlock(2, 4, sizeX-4, sizeY-5, 0)
	line := Drawing.CreateLine(2, 3, 10)
	modal.AddDrawing(textbox)
	modal.AddDrawing(title)
	modal.AddDrawing(status)
	modal.AddDrawing(line)
	modal.AddDrawing(titleField)
	modal.SetPos(x, y)
	modal.SetVisibility(false)
	modal.SetActive(false)
	return &ViewModal{
		Modal:   modal,
		textbox: textbox,
		title:   titleField,
		status:  status,
	}
}

func (v *ViewModal) Open(title, text, status string) {
	v.Modal.SetVisibility(true)
	v.title.SetText(title)
	v.textbox.ClearAll()
	v.textbox.Paste(text)
	v.status.SetText(fmt.Sprintf("(%v)", status))
}

func (v *ViewModal) Change(title, text, status string) {
	v.title.SetText(title)
	v.textbox.ClearAll()
	v.textbox.Paste(text)
	v.status.SetText(fmt.Sprintf("(%v)", status))
}
func (v *ViewModal) Close() {
	v.Modal.SetVisibility(false)
}

func (v *ViewModal) IsOpen() bool {
	return v.Modal.GetVisibility()
}
