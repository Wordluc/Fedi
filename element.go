package main

import (
	"time"

	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
)

type Element struct {
	components *Component.Container
	rectangle *Drawing.Rectangle
	xPos int
	yPos int
	textDrawing *Drawing.TextBlock
	buttons []*Component.Button
	indexButton int
}

func CreateElement(x,y int,width,height int) *Element{
	textElement:=Drawing.CreateTextBlock(2,2,width-1,height-4,10)
	edgeElement:=Drawing.CreateRectangle(1,1,width-2,height)
	drawingContainer:= Drawing.CreateContainer(0,0);
   drawingContainer.AddChild(edgeElement)
	drawingContainer.AddChild(textElement)
	doneButton:=Component.CreateButton(width/2-2,height-3,8,3,"Done")
	doneButton.SetOnHover(func (){
		doneButton.GetVisibleArea().SetColor(Color.Get(Color.White,Color.None))
	})
	doneButton.SetOnLeave(func (){
		doneButton.GetVisibleArea().SetColor(Color.Get(Color.Gray,Color.None))
	})
	doneButton.SetOnClick(func (){
		doneButton.GetVisibleArea().SetColor(Color.Get(Color.Blue,Color.None))
		time.AfterFunc(time.Millisecond*1000, func() {
			doneButton.OnRelease(0,0)
		})
	})
	deleteButton:=Component.CreateButton(width/2-10,height-3,8,3,"Delete")
	deleteButton.SetOnHover(func (){
		deleteButton.GetVisibleArea().SetColor(Color.Get(Color.White,Color.None))
	})
	deleteButton.SetOnLeave(func (){
		deleteButton.GetVisibleArea().SetColor(Color.Get(Color.Gray,Color.None))
	})
	deleteButton.SetOnClick(func (){
		deleteButton.GetVisibleArea().SetColor(Color.Get(Color.Blue,Color.None))
		time.AfterFunc(time.Millisecond*1000, func() {
			deleteButton.OnRelease(0,0)
		})
	})
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
		buttons:[]*Component.Button{doneButton,deleteButton},
		indexButton:0,
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

func (e *Element) SetVisibility(visible bool) {
	e.components.GetGraphics().SetVisibility(visible)
}

func (e *Element) ChangeButton(){
	e.buttons[e.indexButton].OnOut(0,0)
	e.buttons[e.indexButton].OnRelease(0,0)
	if e.indexButton==0{
		e.indexButton=1
	}else if e.indexButton==1{
		e.indexButton=0
	}
	e.buttons[e.indexButton].OnHover(0,0)
}
