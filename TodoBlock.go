package main

import (
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
)

type TodoBlock struct {
	*Drawing.Container
	text       *Drawing.TextBlock
	title      *Drawing.TextBlock
	id         *Drawing.TextField
	date       *Drawing.TextBlock
	mark       *Drawing.TextField
	line       *Drawing.Line
	isSelected bool
}

func CreateTodoBlock(x, y int, xSize int) *TodoBlock {
	container := Drawing.CreateContainer(x, y)
	line := Drawing.CreateLine(x, y, 1)
	line.SetColor(Color.Get(Color.Red, Color.None))
	mark := Drawing.CreateTextField(x+2, y, "")
	title := Drawing.CreateTextBlock(x+4, y, 20, 1, 0)
	title.SetLayer(2)
	id := Drawing.CreateTextField(x+50, y, "")
	text := Drawing.CreateTextBlock(x+6, y+1, xSize-5, 1, 0)
	date := Drawing.CreateTextBlock(x+23, y, 20, 1, 10)
	date.SetLayer(3)
	container.AddDrawings(line, title, text, date, mark, id)
	container.SetLayer(1)
	return &TodoBlock{
		Container:  container,
		text:       text,
		title:      title,
		mark:       mark,
		date:       date,
		id:         id,
		line:       line,
		isSelected: false,
	}
}

type MarkType = string

const (
	Ready      MarkType = "ready"
	WaitingFor MarkType = "waitingFor"
	Progress   MarkType = "progress"
	Done       MarkType = "done"
	Deleted    MarkType = "deleted"
	Archived   MarkType = "archived"
)

func (t *TodoBlock) SetElement(title, text, date, mark, id string) {
	t.id.SetText(id)
	t.text.SetText(text)
	t.title.SetText(title)
	t.date.SetText(date)
	switch mark {
	case Ready:
		t.mark.SetText("⛟")
		t.mark.SetColor(Color.Get(Color.Gray, Color.None))
	case Done:
		t.mark.SetText("✓")
		t.mark.SetColor(Color.Get(Color.Green, Color.None))
	case Progress:
		t.mark.SetText("⛏")
		t.mark.SetColor(Color.Get(Color.Yellow, Color.None))
	case Deleted:
		t.mark.SetText("⛔")
	case Archived:
		t.mark.SetText("X")
		t.mark.SetColor(Color.Get(Color.Yellow, Color.None))
	case WaitingFor:
		t.mark.SetText("⚠")
		t.mark.SetColor(Color.Get(Color.Yellow, Color.None))
	}

}

func (t *TodoBlock) Select() {
	t.isSelected = true
	t.line.SetVisibility(t.isSelected)
}

func (t *TodoBlock) Deselect() {
	t.isSelected = false
	t.line.SetVisibility(t.isSelected)
}
