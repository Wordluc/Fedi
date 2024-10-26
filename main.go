package main

import (
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)

var core *GTUI.Gtui
var carosello Carosello=*CreateCarosello(0,0,3)
func createLabel(text string) Core.IEntity {
	labelList := Drawing.CreateTextField(0, 0)
	labelList.SetText(text)
	bottonLine := Drawing.CreateLine(0, 1, len(text)+1, 0)
	container := Drawing.CreateContainer(0, 0)
	container.AddChild(labelList)
	container.AddChild(bottonLine)
	return container
}

func main() {
	var e error
	core, e = GTUI.NewGtui(loop, &Keyboard.Keyboard{}, &Terminal.Terminal{})
	if e != nil {
		panic(e)
	}

	xSize, ySize := core.Size()
	listZoneXSize := int(float32(xSize) * 0.7)
	listZone := Drawing.CreateRectangle(0, 0, listZoneXSize-2, ySize)
	core.InsertEntity(listZone)
	insertZone := Drawing.CreateRectangle(listZoneXSize, 0, xSize-listZoneXSize, ySize)
	core.InsertEntity(insertZone)
	listLabel := createLabel("To Do")
	listLabel.SetPos(1, 1)
	core.InsertEntity(listLabel)
	editLabel := createLabel("Edit")
	editLabel.SetPos(listZoneXSize+1, 1)
	core.InsertEntity(editLabel)
	listTexts := []string{"1", "2", "3","4","5","6"}
	listElementYSize := int(float32(ySize) * 0.3)
	var elements []*Element = make([]*Element, 3)
	for i := 0; i < len(elements); i++ {
		elements[i] = CreateElement(0, i*listElementYSize+2, listZoneXSize-4, listElementYSize)
		core.InsertComponent(elements[i].GetComponent())
	}
	for i := 0; i < len(listTexts); i++ {
		caroselloEl := &CaroselloElement{
			wakeUpCallBack: func(index int) {
				elements[index%3].components.SetActivity(true)
				elements[index%3].rectangle.SetColor(Color.Get(Color.Blue, Color.None))
			},
			sleepCallBack: func(index int) {
				elements[index%3].components.SetActivity(false)
				elements[index%3].rectangle.SetColor(Color.Get(Color.Gray, Color.None))
			},
			updateCallBack: func(index int) {
				elements[index%3].SetText(listTexts[index])
			},
		}
		carosello.AddElement(caroselloEl)
	}

	firstEdit := true
	TextBox, e := Component.CreateTextBox(listZoneXSize+1, 5, xSize-listZoneXSize-2, ySize-10, core.CreateStreamingCharacter())
	if e != nil {
		panic(e)
	}
	TextBox.Paste("Here you can write your todo")
	TextBox.SetOnClick(func() {
		if firstEdit {
			firstEdit = false
			TextBox.ClearAll()
		}
	})
	TextBox.SetOnHover(func() {
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.Green, Color.None))
	})

	TextBox.SetOnOut(func() {
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
	})

	SendButton := Component.CreateButton(listZoneXSize+1, ySize-5, 8, 3, "Send")
	CancelButton := Component.CreateButton(listZoneXSize+17, ySize-5, 8, 3, "Cancel")
	CancelButton.SetOnClick(func() {
		TextBox.ClearAll()
	})
	core.InsertComponent(TextBox)
	core.InsertComponent(SendButton)
	core.InsertComponent(CancelButton)
	defer core.Start()
}

func loop(keyb Keyboard.IKeyBoard) bool {
	var x, y = core.GetCur()
	if keyb.IsKeySPressed(Keyboard.KeyArrowDown) {
		carosello.NextOrPre(false)
	}
	if keyb.IsKeySPressed(Keyboard.KeyArrowUp) {
		carosello.NextOrPre(true)
	}
	if keyb.IsKeySPressed(Keyboard.KeyEsc) {
		core.EventOn(x, y, func(c Component.IComponent) {
			if c, ok := c.(Component.IWritableComponent); ok {
				c.StopTyping()
			}
		})
	}

	if keyb.IsKeySPressed(Keyboard.KeyEnter) {
		core.Click(x, y)
	}

	if keyb.IsKeySPressed(Keyboard.KeyCtrlV) {
		core.EventOn(x, y, func(c Component.IComponent) {
			if c, ok := c.(*Component.TextBox); ok {
				if c.IsTyping() {
					c.Paste(keyb.GetClickboard())
					core.AllineCursor()
				}
			}
		})
	}

	if keyb.IsKeySPressed(Keyboard.KeyCtrlA) {
		core.EventOn(x, y, func(c Component.IComponent) {
			if c, ok := c.(*Component.TextBox); ok {
				if c.IsTyping() {
					c.SetWrap(!c.GetWrap())
				}
			}
		})
	}

	if keyb.IsKeySPressed(Keyboard.KeyCtrlC) {
		core.EventOn(x, y, func(c Component.IComponent) {
			if c, ok := c.(*Component.TextBox); ok {
				if c.IsTyping() {
					keyb.InsertClickboard(c.Copy())
				}
			}
		})
	}

	if keyb.IsKeySPressed(Keyboard.KeyCtrlQ) {
		return false
	}
	return true
}
