package main

import (
	"Fedi/StateMachine"

	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
)

var core *GTUI.Gtui
var carosello Carosello = *CreateCarosello(0, 0, 3)
var todoBlock []*TodoBlock = make([]*TodoBlock, 3)
var keyb Keyboard.IKeyBoard
var stataMachine *StateMachine.StateMachine
var x, y = 0, 0
var client IApi
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
	keyb = Keyboard.NewKeyboard()
	core, e = GTUI.NewGtui(loop, keyb, &Terminal.Terminal{})
	if e != nil {
		panic(e)
	}
	client,e=CreateNotionClient(".env")
	xSize, ySize := core.Size()
	listZoneXSize := int(float32(xSize) * 0.7)
	todoRect := Drawing.CreateRectangle(0, 0, listZoneXSize-1, ySize)
	todoRect.SetColor(Color.Get(Color.Gray, Color.None))
	core.InsertEntity(todoRect)
	editRect := Drawing.CreateRectangle(listZoneXSize, 0, xSize-listZoneXSize, ySize)
	editRect.SetColor(Color.Get(Color.Gray, Color.None))
	core.InsertEntity(editRect)
	listLabel := createLabel("To Do")
	listLabel.SetPos(1, 1)
	core.InsertEntity(listLabel)
	editLabel := createLabel("Edit")
	editLabel.SetPos(listZoneXSize+1, 1)
	core.InsertEntity(editLabel)
	
	todos,e:=client.GetTodos()
	if e!=nil{
		panic(e)
	}
	listElementYSize := int(float32(ySize) * 0.3)
	for i := 0; i < len(todoBlock); i++ {
		todoBlock[i] = CreateElement(1, i*listElementYSize+3, listZoneXSize-4, listElementYSize)
		core.InsertComponent(todoBlock[i].GetComponent())
	}
	for i := 0; i < len(todos.Todos); i++ {
		caroselloEl := &CaroselloElement{
			wakeUpCallBack: func(todoBlockToUpdate int) {
				todoBlock[todoBlockToUpdate].components.SetActivity(true)
				todoBlock[todoBlockToUpdate].rectangle.SetColor(Color.Get(Color.White, Color.None))
			},
			sleepCallBack: func(todoBlockToUpdate int) {
				todoBlock[todoBlockToUpdate].components.SetActivity(false)
				todoBlock[todoBlockToUpdate].rectangle.SetColor(Color.Get(Color.Gray, Color.None))
			},
			updateCallBack: func(todoBlockToUpdate int) {
				todoBlock[todoBlockToUpdate].SetText(todos.Todos[i].Description)
				todoBlock[todoBlockToUpdate].SetTitle(todos.Todos[i].Name)
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
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.Green, Color.None))
	})
	TextBox.SetOnHover(func() {
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
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

	stataMachine = StateMachine.CreateStateMachine()
	{
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
			} else if keyb.IsKeySPressed(Keyboard.Down) {
				carosello.NextOrPre(false)
			}
			return nil
		})

		caroselloState.SetEntryAction(func() error {
			carosello.updateElement(false)
			return nil
		})

		bottonsCaroselloState := StateMachine.CreateBuilderStateBase("BottonsState")
		bottonsCaroselloState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter) {
				todoBlock[carosello.index%3].GetCurrentBotton().OnClick(0, 0)
			} else if keyb.IsKeySPressed(Keyboard.Left) {
				todoBlock[carosello.index%3].ChangeButton(DeleteBotton)
			} else if keyb.IsKeySPressed(Keyboard.Right) {
				todoBlock[carosello.index%3].ChangeButton(DoneBotton)
			} else if keyb.IsKeySPressed(Keyboard.Up) {
				carosello.NextOrPre(true)
			} else if keyb.IsKeySPressed(Keyboard.Down) {
				carosello.NextOrPre(false)
			}
			return nil
		})
		bottonsCaroselloState.SetEntryAction(func() error {
			todoBlock[carosello.index%3].Active()
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
			carosello.ForEachElements(func(element *CaroselloElement,todoBlockToUpdate int) {
				element.sleepCallBack(todoBlockToUpdate)
			})
			editRect.SetColor(Color.Get(Color.White, Color.None))
			todoRect.SetColor(Color.Get(Color.Gray, Color.None))
			return nil
		})

		textBoxState := StateMachine.CreateBuilderStateBase("TextBoxState")
		textBoxState.SetEntryAction(func() error {
			TextBox.OnHover(0,0)
			return nil
		})
		firstEdit = true
		textBoxState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter) {
				if !TextBox.IsTyping() {
					if firstEdit {
						TextBox.ClearAll()
						firstEdit = false
					}
					core.SetVisibilityCursor(true)
					x, y = TextBox.GetPos()
					x++
					y++
					TextBox.OnClick(0,0)
				}
			}
			return nil
		})
		textBoxState.SetExitAction(func() error {
			TextBox.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
			core.SetVisibilityCursor(false)
			x, y = 0, 0
			TextBox.StopTyping()
			TextBox.OnOut(0, 0)
			return nil
		})

		bottonSendEditState := StateMachine.CreateBuilderStateBase("BottonsEditState")
		bottonSendEditState.SetEntryAction(func() error {
			SendButton.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
			return nil
		})
		bottonSendEditState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter) {
				SendButton.OnClick(0, 0)
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
				CancelButton.OnClick(0, 0)
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
		editPart.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Left)
		}, todoPart)
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Right)
		}, bottonCancelEditState)
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Left)
		}, bottonSendEditState)
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editPart)
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editPart)
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Left)
		}, todoPart)
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

func loop(keyb Keyboard.IKeyBoard) bool {
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
	core.SetCur(x, y)
	stataMachine.Clock()
	core.SetCur(x, y)
	return true
}
