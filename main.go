package main

import "fmt"

import (
	"Fedi/Api"
	"Fedi/StateMachine"
	"slices"
	"strings"
	"time"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)

var carosello Carosello = *CreateCarosello(0, 0, 3)
var todoBlock []*TodoBlock = make([]*TodoBlock, 3)
var stataMachine *StateMachine.StateMachine
var x, y = 0, 0
var client Api.IApi

func createLabel(text string) Core.IEntity {
	labelList := Drawing.CreateTextField(0, 0)
	labelList.SetText(text)
	bottonLine := Drawing.CreateLine(0, 1, len(text)+1, 0)
	container := Drawing.CreateContainer(0, 0)
	container.AddChild(labelList)
	container.AddChild(bottonLine)
	return container
}
func createCarosleloElement(todo Api.Todo) *CaroselloElement {
	return &CaroselloElement{
		wakeUpCallBack: func(todoBlockToUpdate int) {
			todoBlock[todoBlockToUpdate].components.SetActivity(true)
			todoBlock[todoBlockToUpdate].rectangle.SetColor(Color.Get(Color.White, Color.None))
		},
		sleepCallBack: func(todoBlockToUpdate int) {
			todoBlock[todoBlockToUpdate].components.SetActivity(false)
			todoBlock[todoBlockToUpdate].rectangle.SetColor(Color.Get(Color.Gray, Color.None))
		},
		updateCallBack: func(todoBlockToUpdate int) {
			todoBlock[todoBlockToUpdate].SetText(todo.Description)
			todoBlock[todoBlockToUpdate].SetTitle(todo.Name)
		},
	}

}
func main() {
	var e error
	keyb := Keyboard.NewKeyboard()
	core, e := GTUI.NewGtui(loop, keyb, &Terminal.Terminal{})
	if e != nil {
		panic(e)
	}
	client, e = Api.CreateClient(".env")
	if e != nil {
		panic(e)
	}
	xSize, ySize := core.Size()
	listZoneXSize := int(float32(xSize) * 0.7)
	todoRect := Drawing.CreateRectangle(0, 0, listZoneXSize-1, ySize)
	todoRect.SetColor(Color.Get(Color.Gray, Color.None))
	editRect := Drawing.CreateRectangle(listZoneXSize, 0, xSize-listZoneXSize, ySize)
	editRect.SetColor(Color.Get(Color.Gray, Color.None))
	listLabel := createLabel("To Do")
	listLabel.SetPos(1, 1)
	editLabel := createLabel("Edit")
	editLabel.SetPos(listZoneXSize+1, 1)

	todos, e := client.GetTodos()
	if e != nil {
		panic(e)
	}
	listElementYSize := int(float32(ySize) * 0.3)
	for i := 0; i < len(todoBlock); i++ {
		todoBlock[i] = CreateElement(1, i*listElementYSize+3, listZoneXSize-4, listElementYSize)
		core.InsertComponent(todoBlock[i].GetComponent())
	}
	for i := 0; i < len(todos.Todos); i++ {
		carosello.AddElement(createCarosleloElement(todos.Todos[i]))
	}

	numberOfTodoLabel := Drawing.CreateTextField(listZoneXSize-8, 2)
	numberOfTodoLabel.SetText(fmt.Sprint("0/", len(carosello.elements), "  "))
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
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.Green, Color.None))
	})
	TextBox.SetOnHover(func() {
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
	})
	TextBox.SetOnLeave(func() {
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
	})
	SendButton := Component.CreateButton(listZoneXSize+1, ySize-5, 8, 3, "Send")
	SendButton.SetOnClick(func() {
		SendButton.GetVisibleArea().SetColor(Color.Get(Color.Green, Color.None))
		text := TextBox.GetText()
		if strings.TrimFunc(text, func(r rune) bool { return slices.Contains([]rune{' ', '\t', '\n', '\r'}, r) }) == "" {
			return
		}
		TextBox.ClearAll()
		go func() {
			body := Api.Todos{Todos: []Api.Todo{{Description: text, Name: "titolo"}}}
			error := client.PostTodos(body)
			if error != nil {
				panic(error)
			}
			carosello.AddElement(createCarosleloElement(Api.Todo{Description: text, Name: "titolo"}))
			numberOfTodoLabel.SetText(fmt.Sprint(carosello.GetIntex(), "/", len(carosello.elements), " "))
		}()
	})
	SendButton.SetOnRelease(func() {
		SendButton.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
	})
	SendButton.OnRelease()
	CancelButton := Component.CreateButton(listZoneXSize+17, ySize-5, 8, 3, "Cancel")
	CancelButton.SetOnRelease(func() {
		CancelButton.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
	})
	CancelButton.OnRelease()
	CancelButton.SetOnClick(func() {
		CancelButton.GetVisibleArea().SetColor(Color.Get(Color.Green, Color.None))
		TextBox.ClearAll()
		time.AfterFunc(time.Millisecond*1000, func() {
			CancelButton.OnRelease()
		})
	})
	core.InsertComponent(TextBox)
	core.InsertComponent(SendButton)
	core.InsertComponent(CancelButton)

	core.InsertEntity(todoRect)
	core.InsertEntity(editRect)
	core.InsertEntity(listLabel)
	core.InsertEntity(editLabel)
	core.InsertEntity(numberOfTodoLabel)
	stataMachine = StateMachine.CreateStateMachine()
	{ //State machine
		todoPart := StateMachine.CreateBuilderStateBase("todoPart")

		todoPart.SetEntryAction(func() error {
			todoRect.SetColor(Color.Get(Color.White, Color.None))
			editRect.SetColor(Color.Get(Color.Gray, Color.None))
			return nil
		})

		caroselloState := StateMachine.CreateBuilderStateBase("caroselloState")
		caroselloState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Up) {
				carosello.NextOrPre(true)
				numberOfTodoLabel.SetText(fmt.Sprint(carosello.GetIntex(), "/", len(carosello.elements), " "))
			} else if keyb.IsKeySPressed(Keyboard.Down) {
				carosello.NextOrPre(false)
				numberOfTodoLabel.SetText(fmt.Sprint(carosello.GetIntex(), "/", len(carosello.elements), " "))
			} else if keyb.IsKeySPressed(Keyboard.CtrlUp) {
				carosello.SetIndex(0)
			}
			return nil
		})

		caroselloState.SetEntryAction(func() error {
			carosello.UpdateElement(false)
			return nil
		})
		caroselloState.SetExitAction(func() error {
			carosello.SleepAll()
			return nil
		})
		bottonsCaroselloState := StateMachine.CreateBuilderStateBase("BottonsState")
		bottonsCaroselloState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter) {
				todoBlock[carosello.GetSelected()].GetCurrentBotton().OnClick()
			} else if keyb.IsKeySPressed(Keyboard.Left) {
				todoBlock[carosello.GetSelected()].ChangeButton(DeleteBotton)
			} else if keyb.IsKeySPressed(Keyboard.Right) {
				todoBlock[carosello.GetSelected()].ChangeButton(DoneBotton)
			} else if keyb.IsKeySPressed(Keyboard.Up) {
				carosello.NextOrPre(true)
			} else if keyb.IsKeySPressed(Keyboard.Down) {
				carosello.NextOrPre(false)
			}
			return nil
		})
		bottonsCaroselloState.SetEntryAction(func() error {
			todoBlock[carosello.GetSelected()].Active()
			return nil
		})
		bottonsCaroselloState.SetExitAction(func() error {
			for i := 0; i < len(todoBlock); i++ {
				todoBlock[i].Inactive()
			}
			return nil
		})
		editPart := StateMachine.CreateBuilderStateBase("editPart")

		editPart.SetEntryAction(func() error {
			carosello.ForEachElements(func(element *CaroselloElement, todoBlockToUpdate int) {
				element.sleepCallBack(todoBlockToUpdate)
			})
			editRect.SetColor(Color.Get(Color.White, Color.None))
			todoRect.SetColor(Color.Get(Color.Gray, Color.None))
			return nil
		})

		textBoxState := StateMachine.CreateBuilderStateBase("TextBoxState")
		textBoxState.SetEntryAction(func() error {
			TextBox.OnHover()
			return nil
		})
		firstEdit = true
		textBoxState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter) {
				if !TextBox.IsTyping() {
					TextBox.OnRelease()
					if firstEdit {
						TextBox.ClearAll()
						firstEdit = false
					}
					core.SetVisibilityCursor(true)
					x, y = TextBox.GetPos()
					x++
					y++
					TextBox.OnClick()
				}
			}
			if keyb.IsKeySPressed(Keyboard.CtrlC) && TextBox.IsInSelectingMode() {
				keyb.InsertClickboard(TextBox.GetSelectedText())
			}
			if keyb.IsKeySPressed(Keyboard.CtrlV) {
				TextBox.Paste(keyb.GetClickboard())
			}
			if keyb.IsKeySPressed(Keyboard.CtrlA) {
				if TextBox.IsTyping() {
					TextBox.SetWrap(!TextBox.IsInSelectingMode())
				}
			}
			return nil
		})
		textBoxState.SetExitAction(func() error {
			TextBox.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
			core.SetVisibilityCursor(false)
			x, y = 0, 0
			TextBox.OnLeave()
			return nil
		})

		bottonSendEditState := StateMachine.CreateBuilderStateBase("BottonsEditState")
		bottonSendEditState.SetEntryAction(func() error {
			SendButton.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
			return nil
		})
		bottonSendEditState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter) {
				SendButton.OnClick()
			}
			return nil
		})
		bottonSendEditState.SetExitAction(func() error {
			SendButton.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
			return nil
		})

		bottonCancelEditState := StateMachine.CreateBuilderStateBase("BottonsEditState")
		bottonCancelEditState.SetEntryAction(func() error {
			CancelButton.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
			return nil
		})
		bottonCancelEditState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter) {
				CancelButton.OnClick()
			}
			return nil
		})
		bottonCancelEditState.SetExitAction(func() error {
			CancelButton.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
			return nil
		})

		editPart.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Enter)
		}, textBoxState)
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editPart)
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlDown)
		}, bottonSendEditState)
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlLeft)
		}, caroselloState)
		editPart.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Left)
		}, todoPart)
		bottonsCaroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlRight)
		}, textBoxState)
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Right)
		}, bottonCancelEditState)
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Left)
		}, bottonSendEditState)
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlLeft)
		}, caroselloState)
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editPart)
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editPart)
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Left)
		}, todoPart)
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlLeft)
		}, caroselloState)
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Up)
		}, textBoxState)
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Up)
		}, textBoxState)
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editPart)
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Left) && !TextBox.IsTyping()
		}, todoPart)
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Down) && !TextBox.IsTyping()
		}, bottonSendEditState)
		todoPart.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Right)
		}, editPart)
		todoPart.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Enter)
		}, caroselloState)
		caroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, todoPart)
		caroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlRight)
		}, textBoxState)
		caroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Right)
		}, editPart)
		caroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Enter)
		}, bottonsCaroselloState)
		bottonsCaroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc) || keyb.IsKeySPressed(Keyboard.Up) || keyb.IsKeySPressed(Keyboard.Down)
		}, caroselloState)
		stataMachine.AddBuilder(todoPart)
	}

	core.SetVisibilityCursor(false)
	core.Start()
}

func loop(keyb Keyboard.IKeyBoard, core *GTUI.Gtui) bool {
	x, y = core.GetCur()
	if keyb.IsKeySPressed(Keyboard.Left) {
		x--
	} else {
		if keyb.IsKeySPressed(Keyboard.Right) {
			x++
		}
	}
	if keyb.IsKeySPressed(Keyboard.Up) {
		y--
	} else {
		if keyb.IsKeySPressed(Keyboard.Down) {
			y++
		}
	}
	if keyb.IsKeySPressed(Keyboard.CtrlQ) {
		return false
	}
	stataMachine.Clock()
	core.SetCur(x, y)
	return true
}
