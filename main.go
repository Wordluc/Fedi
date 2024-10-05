package main

import (
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)
var core GTUI.Gtui
func createElement(text string,width,height int) Core.IEntity{
	textElement:=Drawing.CreateTextBlock(2,2,width-1,height-4,len(text))
	for i:=range text{
		textElement.Type(rune(text[i]))
	}
	edgeElement:=Drawing.CreateRectangle(1,1,width-2,height)
	container:= Drawing.CreateContainer(0,0);
   container.AddChild(edgeElement)
	container.AddChild(textElement)
	doneButton:=Component.CreateButton(width/2-2,height-3,8,3,"Done")
	deleteButton:=Component.CreateButton(width/2-10,height-3,8,3,"Delete")
	editButton:=Component.CreateButton(width/2+6,height-3,8,3,"Edit")
	container.AddChild(editButton.GetGraphics())
	container.AddChild(deleteButton.GetGraphics())
	container.AddChild(doneButton.GetGraphics())//no, bisogna sistamare il framework

	return container
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
	core,e:=GTUI.NewGtui(loop,&Keyboard.Keyboard{},&Terminal.Terminal{})
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
	for i:=0;i<len(listTexts);i++{
		element:=createElement(listTexts[i],listZoneXSize-2,listElementYSize)
		element.SetPos(0,i*listElementYSize+2)
		core.InsertEntity(element)
	}

	TextBox:=Component.CreateTextBox(listZoneXSize+1,5,xSize-listZoneXSize-2,ySize-10,core.CreateStreamingCharacter())
	TextBox.Paste("Prova di editing")
	SendButton:=Component.CreateButton(listZoneXSize+1,ySize-5,8,3,"Send")
	EditButton:=Component.CreateButton(listZoneXSize+9,ySize-5,8,3,"Edit")
	CancelButton:=Component.CreateButton(listZoneXSize+17,ySize-5,8,3,"Cancel")
	core.InsertComponent(TextBox)
	core.InsertComponent(SendButton)
	core.InsertComponent(EditButton)
	core.InsertComponent(CancelButton)
	core.InsertEntity(TextBox.GetGraphics())

	defer core.Start()
}

func loop(keyboard Keyboard.IKeyBoard) bool {
	if keyboard.IsKeyPressed('c'){
		return false
	}
	return true
}
