package main

import (
	"time"

	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
)
type BottonType int8

const (
	DeleteBotton=iota
	DoneBotton
)
type TodoBlock struct {
	components *Component.Container
	rectangle *Drawing.Rectangle
	xPos int
	yPos int
	textDrawing *Drawing.TextBlock
	buttons []*Component.Button
	currentBottonType BottonType
}

func CreateElement(x,y int,width,height int) *TodoBlock{
	textElement:=Drawing.CreateTextBlock(2,2,width-2,height-4,10)
	edgeElement:=Drawing.CreateRectangle(1,1,width-2,height)
	edgeElement.SetColor(Color.Get(Color.Gray,Color.None))
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
	containerComponent.AddDrawing(drawingContainer)
	containerComponent.SetPos(x,y)
	return &TodoBlock{
		components:containerComponent,
		rectangle:edgeElement,
		xPos:x,
		yPos:y,
		textDrawing:textElement,
		buttons:[]*Component.Button{deleteButton,doneButton},
	}
}

func (e *TodoBlock) SetPos(x,y int){
	e.components.SetPos(x,y)
	e.xPos=x
	e.yPos=y
}

func (e *TodoBlock) GetPos() (int,int){
	return e.xPos,e.yPos
}
func (e *TodoBlock) GetComponent() *Component.Container{
	return e.components
}
func (e *TodoBlock) SetText(text string){
	e.textDrawing.ClearAll()
	for i:=range text{
		e.textDrawing.Type(rune(text[i]))
	}
}

func (e *TodoBlock) SetVisibility(visible bool) {
	e.components.GetGraphics().SetVisibility(visible)
}

func (e *TodoBlock) ChangeButton(bottontype BottonType){
	if bottontype!=e.currentBottonType{
		e.buttons[e.currentBottonType].OnOut(0,0)
		e.buttons[e.currentBottonType].OnRelease(0,0)
		e.currentBottonType=bottontype
	}
	e.buttons[e.currentBottonType].OnHover(0,0)
}

func (e *TodoBlock) GetCurrentBotton() *Component.Button{
	return e.buttons[e.currentBottonType]
}

func (e *TodoBlock) Active(){
	e.ChangeButton(DeleteBotton)
}
func (e *TodoBlock) Inactive(){
	for i:=range e.buttons{
		e.buttons[i].OnOut(0,0)
		e.buttons[i].OnRelease(0,0)
	}
}