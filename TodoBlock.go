package main

import (
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
)

type TodoBlock struct {
	*Drawing.Container
	text       *Drawing.TextBlock
	title      *Drawing.TextBlock
	date       *Drawing.TextBlock
	line       *Drawing.Line
	isSelected bool
}

func CreateTodoBlock(x, y int, xSize int) *TodoBlock {
	container := Drawing.CreateContainer(x, y)
	cursor := Drawing.CreateLine(x, y, 1)
	cursor.SetColor(Color.Get(Color.Red, Color.None))
	title := Drawing.CreateTextBlock(x+2, y, 20, 1, 0)
	text := Drawing.CreateTextBlock(x+4, y+1, xSize, 1, 0)
	date := Drawing.CreateTextBlock(x+23, y, 20, 1, 10)
	container.AddDrawings(cursor, title, text, date)
	container.SetLayer(1)
	return &TodoBlock{
		Container:  container,
		text:       text,
		title:      title,
		date:       date,
		line:       cursor,
		isSelected: false,
	}
}

func (t *TodoBlock) SetElement(title, text, date string) {
	t.text.SetText(text)
	t.title.SetText(title)
	t.date.SetText(date)
}

func (t *TodoBlock) Select() {
	t.isSelected = true
	t.line.SetVisibility(t.isSelected)
}

func (t *TodoBlock) Deselect() {
	t.isSelected = false
	t.line.SetVisibility(t.isSelected)
}
