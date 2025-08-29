package main

import (
	"strings"

	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
)

type TodoBlock struct {
	*Drawing.Container
	text        *Drawing.TextBlock
	title       *Drawing.TextBlock
	id          *Drawing.TextField
	date        *Drawing.TextBlock
	mark        *Drawing.TextBlock
	line        *Drawing.Line
	iconTextBig *Drawing.TextBlock
	isSelected  bool
}

const sizeTitle = 40

func CreateTodoBlock(x, y int, xSize int) *TodoBlock {
	container := Drawing.CreateContainer()
	line := Drawing.CreateLine(0, 0, 1)
	line.SetColor(Color.Get(Color.Red, Color.None))
	mark := Drawing.CreateTextBlock(2, 0, 1, 1, 1)
	title := Drawing.CreateTextBlock(4, 0, sizeTitle, 1, 0)
	id := Drawing.CreateTextField(sizeTitle+30, 0, "")
	text := Drawing.CreateTextBlock(5, 1, xSize-5, 2, 0)
	horizontalLine := Drawing.CreateLine(4, 1, 2)
	horizontalLine.SetAngle(90)
	iconTextBig := Drawing.CreateTextBlock(4, 2, 1, 1, 0)
	iconTextBig.SetText("⇣")
	iconTextBig.SetVisibility(false)
	date := Drawing.CreateTextBlock(sizeTitle+4, 0, 20, 1, 10)
	container.AddDrawings(line, title, text, date, mark, id, horizontalLine, iconTextBig)
	container.SetPos(x, y)
	return &TodoBlock{
		Container:   container,
		text:        text,
		title:       title,
		mark:        mark,
		date:        date,
		id:          id,
		line:        line,
		isSelected:  false,
		iconTextBig: iconTextBig,
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
	text = strings.ReplaceAll(text, "/n", "\n")
	t.text.SetCursor_Relative(-3, -3)
	t.text.SetText(text)
	t.iconTextBig.SetVisibility(len(strings.Split(text, "\n")) > 2)
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
		t.mark.SetColor(Color.Get(Color.Cyan, Color.None))
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
