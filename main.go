package main

import (
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)
var core *GTUI.Gtui
func createElement(text string,width,height int) *Component.Container{
	textElement:=Drawing.CreateTextBlock(2,2,width-1,height-4,len(text))
	for i:=range text{
		textElement.Type(rune(text[i]))
	}
	edgeElement:=Drawing.CreateRectangle(1,1,width-2,height)
	drawingContainer:= Drawing.CreateContainer(0,0);
   drawingContainer.AddChild(edgeElement)
	drawingContainer.AddChild(textElement)

	doneButton:=Component.CreateButton(width/2-2,height-3,8,3,"Done")
	deleteButton:=Component.CreateButton(width/2-10,height-3,8,3,"Delete")
	editButton:=Component.CreateButton(width/2+6,height-3,8,3,"Edit")
	containerComponent:=Component.CreateContainer(0,0)
	containerComponent.AddComponent(doneButton)
	containerComponent.AddComponent(deleteButton)
	containerComponent.AddComponent(editButton)
	containerComponent.AddDrawing(*drawingContainer)

	return containerComponent
}
func createLabel(text string) Core.IEntity{
	labelList:=Drawing.CreateTextField(0,0)
	labelList.SetText(text)
	bottonLine:=Drawing.CreateLine(0,1,len(text)+1,0)
   container:=Drawing.CreateContainer(0,0)
	container.AddChild(labelList)
	container.AddChild(bottonLine)
	return container
}

func main() {
	var e error
	core,e=GTUI.NewGtui(loop,&Keyboard.Keyboard{},&Terminal.Terminal{})
	if e!=nil{
		panic(e)
	}

	xSize,ySize:=core.Size()
	listZoneXSize:=int(float32(xSize)*0.7)
	listZone:=Drawing.CreateRectangle(0,0,listZoneXSize-2,ySize)
	core.InsertEntity(listZone)
	insertZone:=Drawing.CreateRectangle(listZoneXSize,0,xSize-listZoneXSize,ySize)
	core.InsertEntity(insertZone)
	listLabel:=createLabel("To Do")
	listLabel.SetPos(1,1)
	core.InsertEntity(listLabel)
	editLabel:=createLabel("Edit")
	editLabel.SetPos(listZoneXSize+1,1)
	core.InsertEntity(editLabel)
	listTexts:=[]string{"Prova","fai questa cosa","Esplodi"}
	listElementYSize:=int(float32(ySize)*0.3)
   var elements []*Component.Container=make([]*Component.Container,len(listTexts))
	for i:=0;i<len(listTexts);i++{
		elements[i]=createElement(listTexts[i],listZoneXSize-2,listElementYSize)
		elements[i].SetActivity(false)
		elements[i].SetPos(0,i*listElementYSize+2)
		core.InsertComponent(elements[i])
	}
	x,y:=elements[0].GetGraphics().GetPos()
	elements[0].SetActivity(true)
	core.SetCur(x+2,y+2)
	firstEdit:=true
	TextBox:=Component.CreateTextBox(listZoneXSize+1,5,xSize-listZoneXSize-2,ySize-10,core.CreateStreamingCharacter())
	TextBox.Paste("Here you can write your todo")
	TextBox.SetOnClick(func() {
		if firstEdit {
			firstEdit=false
			TextBox.ClearAll()
		}
	})
	SendButton:=Component.CreateButton(listZoneXSize+1,ySize-5,8,3,"Send")
	CancelButton:=Component.CreateButton(listZoneXSize+17,ySize-5,8,3,"Cancel")
	CancelButton.SetOnClick(func() {
		TextBox.ClearAll()
	})
	core.InsertComponent(TextBox)
	core.InsertComponent(SendButton)
	core.InsertComponent(CancelButton)
	core.InsertEntity(TextBox.GetGraphics())
	defer core.Start()
}

func loop(keyb Keyboard.IKeyBoard) bool {
	var x, y = core.GetCur()
	if keyb.IsKeySPressed(Keyboard.KeyArrowDown) {
		y++
	}
	if keyb.IsKeySPressed(Keyboard.KeyArrowUp) {
		y--
	}
	if keyb.IsKeySPressed(Keyboard.KeyArrowRight) {
		x++
	}
	if keyb.IsKeySPressed(Keyboard.KeyArrowLeft) {
		x--
	}

	if keyb.IsKeySPressed(Keyboard.KeyCtrlS) {
		core.RefreshComponents()
	}

	if e:=core.SetCur(x,y);e!=nil{
		panic(e)
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

