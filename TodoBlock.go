package main

import (
	"Fedi/Api"
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
	titleDrawing *Drawing.TextBlock
	lineTitle *Drawing.Line
	buttons []*Component.Button
	currentBottonType BottonType
	currentTodo *Api.Todo
}

func CreateElement(x,y int,width,height int,toDelete func ()) *TodoBlock{
	title:=Drawing.CreateTextBlock(2,2,width-5,3,10)
	line:=Drawing.CreateLine(2,3,3,0)
	line.SetVisibility(false)
	textElement:=Drawing.CreateTextBlock(3,4,width-5,height-4,10)
	edgeElement:=Drawing.CreateRectangle(1,1,width-2,height)
	edgeElement.SetColor(Color.Get(Color.Gray,Color.None))
	drawingContainer:= Drawing.CreateContainer(0,0);
   drawingContainer.AddChild(edgeElement)
	drawingContainer.AddChild(textElement)
	drawingContainer.AddChild(title)
	drawingContainer.AddChild(line)
	doneButton:=Component.CreateButton(width/2-2,height-3,8,3,"Done")
	doneButton.SetOnHover(func (){
		doneButton.GetVisibleArea().SetBorderColor(Color.Get(Color.White,Color.None))
	})
	doneButton.SetOnLeave(func (){
		doneButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Gray,Color.None))
	})
	doneButton.SetOnRelease(func (){
		doneButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Gray,Color.None))
	})
	doneButton.SetOnClick(func (){
		doneButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Blue,Color.None))
	})
	deleteButton:=Component.CreateButton(width/2-10,height-3,8,3,"Delete")
	deleteButton.SetOnHover(func (){
		deleteButton.GetVisibleArea().SetBorderColor(Color.Get(Color.White,Color.None))
	})
	deleteButton.SetOnLeave(func (){
		deleteButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Gray,Color.None))
	})
	deleteButton.SetOnRelease(func (){
		deleteButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Gray,Color.None))
	})
	deleteButton.SetOnClick(func (){
		deleteButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Blue,Color.None))
			go func (){
				toDelete()
			}()
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
		titleDrawing:title,
		lineTitle:line,
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
func (e *TodoBlock) setText(text string){
	e.textDrawing.ClearAll()
	for i:=range text{
		e.textDrawing.Type(rune(text[i]))
	}
}
func (e *TodoBlock) setTitle(text string){
	e.titleDrawing.SetText(text)
	e.lineTitle.SetVisibility(true)
}

func (e *TodoBlock) Clean() {
	e.setText("")
	e.setTitle("")
	e.lineTitle.SetVisibility(false)
}

func (e *TodoBlock) ChangeButton(bottontype BottonType){
	if bottontype!=e.currentBottonType{
		e.buttons[e.currentBottonType].OnRelease()
		e.currentBottonType=bottontype
	}
	e.buttons[e.currentBottonType].OnHover()
}

func (e *TodoBlock) GetCurrentBotton() *Component.Button{
	return e.buttons[e.currentBottonType]
}
func (e *TodoBlock) SetCurrentTodo(todo *Api.Todo){
	e.currentTodo=todo
	e.setTitle(todo.Name)
	e.setText(todo.Description)
}
func (e *TodoBlock) GetTodo() *Api.Todo{
	return e.currentTodo
}
func (e *TodoBlock) Active(){//todo: da togliere
	e.ChangeButton(DeleteBotton)
}
func (e *TodoBlock) ReleaseAll(){
	for i:=range e.buttons{
		e.buttons[i].OnRelease()
		e.buttons[i].OnRelease()
	}
}
