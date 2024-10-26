package main

import (
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
)

type Element struct {
	components *Component.Container
	rectangle *Drawing.Rectangle
	xPos int
	yPos int
	textDrawing *Drawing.TextBlock
	doneButton *Component.Button
	deleteButton *Component.Button
}

func CreateElement(x,y int,width,height int) *Element{
	textElement:=Drawing.CreateTextBlock(2,2,width-1,height-4,10)
	edgeElement:=Drawing.CreateRectangle(1,1,width-2,height)
	drawingContainer:= Drawing.CreateContainer(0,0);
   drawingContainer.AddChild(edgeElement)
	drawingContainer.AddChild(textElement)
	doneButton:=Component.CreateButton(width/2-2,height-3,8,3,"Done")
	deleteButton:=Component.CreateButton(width/2-10,height-3,8,3,"Delete")
	containerComponent:=Component.CreateContainer(0,0)
	containerComponent.AddComponent(doneButton)
	containerComponent.AddComponent(deleteButton)
	containerComponent.AddDrawing(*drawingContainer)
	containerComponent.SetPos(x,y)
	return &Element{
		components:containerComponent,
		rectangle:edgeElement,
		xPos:x,
		yPos:y,
		textDrawing:textElement,
		doneButton:doneButton,
		deleteButton:deleteButton,
	}
}

func (e *Element) SetPos(x,y int){
	e.components.SetPos(x,y)
	e.xPos=x
	e.yPos=y
}

func (e *Element) GetPos() (int,int){
	return e.xPos,e.yPos
}
func (e *Element) GetComponent() *Component.Container{
	return e.components
}
func (e *Element) SetText(text string){
	e.textDrawing.ClearAll()
	for i:=range text{
		e.textDrawing.Type(rune(text[i]))
	}
}

func (e *Element) SetCallbackOnDone(callback func ()) {
	e.doneButton.SetOnClick(callback)
}

func (e *Element) SetCallbackOnDelete(callback func ()) {
	e.deleteButton.SetOnClick(callback)
}

func (e *Element) SetVisibility(visible bool) {
	e.components.GetGraphics().SetVisibility(visible)
}
