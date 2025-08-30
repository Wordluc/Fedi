package main

import (
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
)

func CreateTutorialModal(w, h int) *Component.Modal {
	const sizeX, sizeY = 81, 43
	modal := Component.CreateModal(sizeX, sizeY)
	x := w/2 - sizeX/2
	y := h/2 - sizeY/2
	title := Drawing.CreateTextField(1, 1, "Tutorial commands")
	commands := [][]string{
		{"Ctrl-q", "Close the app"},
		{"j", "Scroll down"},
		{"k", "Scroll up"},
		{"Ctrl-s", "Open the editing modal and save a new todo, or save changes to an existing one"},
		{"Ctrl-e", "Open the editing modal to edit a todo"},
		{"Ctrl-d", "Set todo as Done"},
		{"Ctrl-x", "Set todo as Deleted"},
		{"Ctrl-a", "Set todo as Archived"},
		{"Ctrl-w", "Set todo as WaitingFor"},
		{"Ctrl-v", "Open open the view modal to see details about the todo"},
		{"(Modal)Ctrl-k", "Focus the title field"},
		{"(Modal)Ctrl-j", "Focus the text field"},
		{"(Modal)Ctrl-q", "Close the modal"},
	}
	const distance = 3
	for c := range commands {
		title := Drawing.CreateTextField(1, 3+distance*c, commands[c][0])
		title.SetColor(Color.Get(Color.Red, Color.None))
		text := Drawing.CreateTextField(2, 4+distance*c, commands[c][1])
		modal.AddDrawing(title, text)
	}
	modal.AddDrawing(title)
	modal.SetPos(x, y)
	return modal
}
