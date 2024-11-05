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
var elements []*Element = make([]*Element, 3)
var keyb Keyboard.IKeyBoard
var stataMachine *StateMachine.StateMachine

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
	core, e = GTUI.NewGtui(loop,keyb , &Terminal.Terminal{})
	if e != nil {
		panic(e)
	}

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
	listTexts := []string{"1", "2", "3", "4", "5", "6"}
	listElementYSize := int(float32(ySize) * 0.3)
	for i := 0; i < len(elements); i++ {
		elements[i] = CreateElement(1, i*listElementYSize+2, listZoneXSize-4, listElementYSize)
		core.InsertComponent(elements[i].GetComponent())
	}
	for i := 0; i < len(listTexts); i++ {
		caroselloEl := &CaroselloElement{
			index: i,
			wakeUpCallBack: func() {
				elements[i%3].components.SetActivity(true)
				elements[i%3].rectangle.SetColor(Color.Get(Color.White, Color.None))
			},
			sleepCallBack: func() {
				elements[i%3].components.SetActivity(false)
				elements[i%3].rectangle.SetColor(Color.Get(Color.Gray, Color.None))
			},
			updateCallBack: func() {
				elements[i%3].SetText(listTexts[i])
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
			}else
			if keyb.IsKeySPressed(Keyboard.Down) {
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
				elements[carosello.index%3].buttons[elements[carosello.index%3].indexButton].OnClick(0,0)
			}else
			if keyb.IsKeySPressed(Keyboard.Left) || keyb.IsKeySPressed(Keyboard.Right)  {
				elements[carosello.index%3].ChangeButton()
			}
			return nil
		})
		bottonsCaroselloState.SetEntryAction(func() error {
			elements[carosello.index%3].buttons[0].OnHover(0,0)
			return nil
		})
		bottonsCaroselloState.SetExitAction(func() error {
			for i := 0; i < len(elements[carosello.index%3].buttons); i++ {
				elements[carosello.index%3].buttons[i].OnOut(0,0)
			}
			return nil
		})
		editPart := StateMachine.CreateBuilderStateBase("editPart")

		editPart.SetEntryAction(func() error {
			carosello.ForEachElements(func(element *CaroselloElement) {
				element.sleepCallBack()
			})
			editRect.SetColor(Color.Get(Color.White, Color.None))
			todoRect.SetColor(Color.Get(Color.Gray, Color.None))
			return nil
		})

		textBoxState := StateMachine.CreateBuilderStateBase("TextBoxState")
		textBoxState.SetEntryAction(func() error {
			TextBox.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
			return nil
		})
		textBoxState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter) {
				TextBox.ClearAll()
				TextBox.StartTyping()
			}
			return nil
		})
		textBoxState.SetExitAction(func() error {
			TextBox.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
			TextBox.StopTyping()
			return nil
		})
		
		bottonSendEditState:=StateMachine.CreateBuilderStateBase("BottonsEditState")
		bottonSendEditState.SetEntryAction(func() error {
			SendButton.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
			return nil
		})
		bottonSendEditState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter){
				SendButton.OnClick(0,0)
			}
			return nil
		})
		bottonSendEditState.SetExitAction(func () error {
			SendButton.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
			return nil
		})

		bottonCancelEditState:=StateMachine.CreateBuilderStateBase("BottonsEditState")
		bottonCancelEditState.SetEntryAction(func() error {
			CancelButton.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
			return nil
		})
		bottonCancelEditState.SetActionDo(func() error {
			if keyb.IsKeySPressed(Keyboard.Enter){
				CancelButton.OnClick(0,0)
			}
			return nil
		})
		bottonCancelEditState.SetExitAction(func () error {
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
			return keyb.IsKeySPressed(Keyboard.Down)
		},bottonSendEditState)
		editPart.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Left)
		},todoPart)
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
		textBoxState.AddBranch(func() bool {
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
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, caroselloState)

		stataMachine.AddBuilder(todoPart)
	}

	defer core.Start()
}

func loop(keyb Keyboard.IKeyBoard) bool {

	if keyb.IsKeySPressed(Keyboard.CtrlQ) {
		return false
	}
	stataMachine.Clock()
   core.IRefreshAll()
	return true
}
