package main

import (
	"Fedi/Api"
	"Fedi/StateMachine"
	"fmt"
	"github.com/Wordluc/GTUI"
	"github.com/Wordluc/GTUI/Core/Component"
	"github.com/Wordluc/GTUI/Core/Drawing"
	"github.com/Wordluc/GTUI/Core/EventManager"
	"github.com/Wordluc/GTUI/Core/Utils/Color"
	"github.com/Wordluc/GTUI/Keyboard"
	"github.com/Wordluc/GTUI/Terminal"
	"slices"
	"strings"
	"time"
)

var carosello *Carosello
var todoBlock []*TodoBlock = make([]*TodoBlock, 3)
var stateMachine *StateMachine.StateMachine
var x, y = 0, 0
var client Api.IApi
var isInDoneState = false

const ClockEvent EventManager.EventType = 10

var numberOfTodoLabel *Drawing.TextBlock
var todos *Api.Todos
var textWarning *Drawing.TextBlock

func createLabel(text string) *Drawing.Container {
	labelList := Drawing.CreateTextField(0, 0, text)
	labelList.SetLayer(2)
	bottonLine := Drawing.CreateLine(0, 1, len(text)+1)
	bottonLine.SetLayer(2)

	container := Drawing.CreateContainer(0, 0)
	container.AddDrawings(labelList, bottonLine)
	return container
}
func createCarosleloElement(todo Api.Todo) *CaroselloElement {
	return &CaroselloElement{
		wakeUpCallBack: func(todoBlockToUpdate int) {
			GTUI.Logf("rect:%v", int(todoBlock[todoBlockToUpdate].rectangle.GetLayer()))
			GTUI.Logf("text:%v", int(todoBlock[todoBlockToUpdate].textDrawing.GetLayer()))
			todoBlock[todoBlockToUpdate].rectangle.SetColor(Color.Get(Color.White, Color.None))
		},
		sleepCallBack: func(todoBlockToUpdate int) {
			todoBlock[todoBlockToUpdate].rectangle.SetColor(Color.Get(Color.Gray, Color.None))
		},
		updateCallBack: func(todoBlockToUpdate int) {
			todoBlock[todoBlockToUpdate].SetCurrentTodo(&todo)
		},
	}
}

func refreshCarosello(carosello **Carosello, allltodos *Api.Todos, setWakeup bool) {
	var filteFor string
	if isInDoneState {
		filteFor = "Done"
	} else {
		filteFor = "Todo"
	}
	todosToShow := filterTodo(allltodos, filteFor)
	newCarosello := CreateCarosello(0, 0, 3)
	for i := 0; i < len(todosToShow.Todos); i++ {
		newCarosello.AddElement(createCarosleloElement(todosToShow.Todos[i]))
	}
	newCarosello.UpdateElementState(true, setWakeup)
	*carosello = newCarosello
	numberOfTodoLabel.SetText(fmt.Sprint("0/", len(newCarosello.elements), "  "))
}

func filterTodo(todos *Api.Todos, on string) *Api.Todos {
	var filteredTodos *Api.Todos = &Api.Todos{}
	for _, todo := range todos.Todos {
		if todo.Status == on {
			filteredTodos.Todos = append(filteredTodos.Todos, todo)
		}
	}
	return filteredTodos
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
	listLabel := createLabel("List of")
	listLabel.SetPos(1, 1)
	editLabel := createLabel("Edit")
	editLabel.SetPos(listZoneXSize+1, 1)
	typeTodoLavel := Drawing.CreateTextBlock(13, 1, 5, 3, 10)
	typeTodoLavel.SetLayer(2)
	typeTodoLavel.SetText("To do")
	typeTodoLavel.SetColor(Color.Get(Color.Red, Color.None))
	core.AddDrawing(typeTodoLavel)

	todos, e = client.GetTodos()
	if e != nil {
		panic(e)
	}
	numberOfTodoLabel = Drawing.CreateTextBlock(listZoneXSize-8, 2, 6, 1, 10)
	numberOfTodoLabel.SetLayer(2)

	listElementYSize := int(float32(ySize) * 0.3)
	for i := 0; i < len(todoBlock); i++ {
		todoBlock[i] = CreateElement(1, i*listElementYSize+3, listZoneXSize-4, listElementYSize, func() { //da ottimizzare
			if carosello.GetElementsNumber() == 0 {
				return
			}
			currentTodo := todoBlock[i].GetTodo()
			if e := client.Delete(*currentTodo); e != nil {
			}
			for i := 0; i < len(todoBlock); i++ {
				todoBlock[i].Clean()
			}
			todos, e = client.GetTodos()
			if e != nil {
				panic(e)
			}
			refreshCarosello(&carosello, todos, true)
			EventManager.Call(ClockEvent, nil)
			EventManager.Call(EventManager.Refresh, []any{todoBlock[i].GetComponent()})
		}, func() {
			if e := client.SetAsDone(*todoBlock[i].GetTodo()); e != nil {
				panic(e)
			}
			todos, e = client.GetTodos()
			if e != nil {
				panic(e)
			}
			refreshCarosello(&carosello, todos, true)
			EventManager.Call(EventManager.Refresh, []any{todoBlock[i].GetComponent()})
		})
		core.AddContainer(todoBlock[i].GetComponent())
	}
	refreshCarosello(&carosello, todos, false)

	TextBox, e := Component.CreateTextBox(listZoneXSize+1, 10, xSize-listZoneXSize-2, ySize-15, core.CreateStreamingCharacter())
	LabelBox := Drawing.CreateTextField(listZoneXSize+1, 9, "Description")
	if e != nil {
		panic(e)
	}
	TextBox.SetOnClick(func() {
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.Green, Color.None))
	})
	TextBox.SetOnHover(func() {
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
	})
	TextBox.SetOnLeave(func() {
		TextBox.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
	})
	LabelTitle := Drawing.CreateTextField(listZoneXSize+1, 4, "Title")

	TitleBox, e := Component.CreateTextBox(listZoneXSize+1, 5, xSize-listZoneXSize-2, 3, core.CreateStreamingCharacter())

	TitleBox.SetOnClick(func() {
		TitleBox.GetVisibleArea().SetColor(Color.Get(Color.Green, Color.None))
	})
	TitleBox.SetOnHover(func() {
		TitleBox.GetVisibleArea().SetColor(Color.Get(Color.White, Color.None))
	})
	TitleBox.SetOnLeave(func() {
		TitleBox.GetVisibleArea().SetColor(Color.Get(Color.Gray, Color.None))
	})

	SendButton := Component.CreateButton(listZoneXSize+1, ySize-5, 8, 3, "Send")
	SendButton.SetOnClick(func() {
		SendButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Green, Color.None))
		title := TitleBox.GetText()
		description := TextBox.GetText()
		if strings.TrimFunc(description, func(r rune) bool { return slices.Contains([]rune{' ', '\t', '\n', '\r'}, r) }) == "" {
			return
		}
		TextBox.ClearAll()
		TitleBox.ClearAll()
		go func() {
			body := Api.Todos{Todos: []Api.Todo{{Description: description, Name: title}}}
			response, error := client.PostTodos(body)
			if error != nil {
				panic(error)
			}
			todos.Todos = append(todos.Todos, Api.Todo{Description: description, Name: title, Id: response.Id, Status: "Todo"})
			refreshCarosello(&carosello, todos, true)
			carosello.UpdateElementState(true, false)
			EventManager.Call(EventManager.Refresh, nil)
		}()
	})
	SendButton.SetOnRelease(func() {
		SendButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Gray, Color.None))
	})
	SendButton.OnRelease()
	CancelButton := Component.CreateButton(listZoneXSize+17, ySize-5, 8, 3, "Cancel")
	CancelButton.SetOnRelease(func() {
		CancelButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Gray, Color.None))
	})
	CancelButton.OnRelease()
	CancelButton.SetOnClick(func() {
		CancelButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Green, Color.None))
		TextBox.ClearAll()
		time.AfterFunc(time.Millisecond*1000, func() {
			CancelButton.OnRelease()
		})
	})
	core.AddComponent(TextBox)
	core.AddComponent(SendButton)
	core.AddComponent(CancelButton)
	core.AddComponent(TitleBox)

	core.AddDrawing(todoRect)
	core.AddDrawing(editRect)
	core.AddContainer(listLabel)
	core.AddContainer(editLabel)
	core.AddDrawing(numberOfTodoLabel)
	core.AddDrawing(LabelBox)
	core.AddDrawing(LabelTitle)
	stateMachine = StateMachine.CreateStateMachine()
	todoState := StateMachine.CreateBuilderStateComposite("todoPart")

	todoState.SetEntryAction(func() error {
		carosello.ForEachElements(func(element *CaroselloElement, todoBlockToUpdate int) {
			element.sleepCallBack(todoBlockToUpdate)
		})
		carosello.ForSelectedElement(func(element *CaroselloElement, todoBlockToUpdate int) {
			element.wakeUpCallBack(todoBlockToUpdate)
		})
		todoRect.SetColor(Color.Get(Color.White, Color.None))
		return nil
	})
	todoState.SetExitAction(func() error {
		todoRect.SetColor(Color.Get(Color.Gray, Color.None))
		return nil
	})
	caroselloState := StateMachine.CreateBuilderStateBase("caroselloState")
	caroselloState.SetActionDo(func() error {
		if keyb.IsKeyPressed('k') {
			carosello.NextOrPre(true)
		} else if keyb.IsKeyPressed('j') {
			carosello.NextOrPre(false)
		} else if keyb.IsKeySPressed(Keyboard.CtrlK) {
			carosello.SetIndex(0)
		} else if keyb.IsKeySPressed(Keyboard.CtrlH) {
			carosello.SetIndex(len(carosello.elements) - 1)
		}
		numberOfTodoLabel.SetText(fmt.Sprint(carosello.GetIntex(), "/", len(carosello.elements), " "))
		return nil
	})

	showTodoState := StateMachine.CreateBuilderStateBase("ShowTodoState")
	showTodoState.SetEntryAction(func() error {
		isInDoneState = false
		keyb.CleanKeyboardState()
		for i := 0; i < len(todoBlock); i++ {
			todoBlock[i].Clean()
		}
		todos, _ := client.GetTodos()
		refreshCarosello(&carosello, todos, false)
		for i := 0; i < len(todoBlock); i++ {
			for _, b := range todoBlock[i].buttons {
				b.GetGraphic().SetVisibility(true)
			}
		}
		numberOfTodoLabel.SetText(fmt.Sprint("0/", len(carosello.elements), "  "))
		EventManager.Call(EventManager.Refresh, []any{todoBlock[carosello.GetSelected()].GetComponent()})
		return nil
	})
	showDoneState := StateMachine.CreateBuilderStateBase("ShowDoneState")
	showDoneState.SetEntryAction(func() error {
		isInDoneState = true
		keyb.CleanKeyboardState()
		for i := 0; i < len(todoBlock); i++ {
			todoBlock[i].Clean()
		}
		todos, _ = client.GetTodos()
		refreshCarosello(&carosello, todos, false)
		for i := 0; i < len(todoBlock); i++ {
			for _, b := range todoBlock[i].buttons {
				b.GetGraphic().SetVisibility(false)
			}
		}
		typeTodoLavel.SetText("Done")
		EventManager.Call(ClockEvent, []any{todoBlock[carosello.GetSelected()].GetComponent()})
		EventManager.Call(EventManager.Refresh, []any{todoBlock[carosello.GetSelected()].GetComponent()})
		numberOfTodoLabel.SetText(fmt.Sprint("0/", len(carosello.elements), "  "))
		return nil
	})
	caroselloState.SetEntryAction(func() error {
		carosello.UpdateElementState(false, true)
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
		} else if keyb.IsKeyPressed('h') {
			todoBlock[carosello.GetSelected()].ChangeButton(DeleteBotton)
		} else if keyb.IsKeyPressed('l') {
			todoBlock[carosello.GetSelected()].ChangeButton(DoneBotton)
		} else if keyb.IsKeyPressed('k') {
			carosello.NextOrPre(true)
		} else if keyb.IsKeyPressed('j') {
			carosello.NextOrPre(false)
		}
		return nil
	})
	bottonsCaroselloState.SetEntryAction(func() error {
		todoBlock[carosello.GetSelected()].Active()

		carosello.ForSelectedElement(func(element *CaroselloElement, todoBlockToUpdate int) {
			element.sleepCallBack(todoBlockToUpdate)
		})
		//		carosello.SleepAll()
		return nil
	})
	bottonsCaroselloState.SetExitAction(func() error {
		for i := 0; i < len(todoBlock); i++ {
			todoBlock[i].ReleaseAll()
		}
		return nil
	})
	editState := StateMachine.CreateBuilderStateComposite("editPart")

	editState.SetEntryAction(func() error {
		editRect.SetColor(Color.Get(Color.White, Color.None))
		return nil
	})
	editState.SetExitAction(func() error {
		editRect.SetColor(Color.Get(Color.Gray, Color.None))
		return nil
	})
	textBoxState := StateMachine.CreateBuilderStateBase("TextBoxState")
	textBoxState.SetEntryAction(func() error {
		TextBox.OnHover()
		return nil
	})
	textBoxState.SetActionDo(func() error {
		if keyb.IsKeySPressed(Keyboard.Enter) {
			if !TextBox.IsTyping() {
				x, y = TextBox.GetPos()
				x++
				y++
				core.SetCur(x, y)
				core.SetVisibilityCursor(true)
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
		core.SetVisibilityCursor(false)
		TextBox.OnLeave()
		return nil
	})
	titleBoxState := StateMachine.CreateBuilderStateBase("TextBoxState")
	titleBoxState.SetEntryAction(func() error {
		TitleBox.OnHover()
		return nil
	})
	titleBoxState.SetActionDo(func() error {
		if keyb.IsKeySPressed(Keyboard.Enter) {
			if !TitleBox.IsTyping() {
				x, y = TitleBox.GetPos()
				x++
				y++
				core.SetCur(x, y)
				core.SetVisibilityCursor(true)
				TitleBox.OnClick()
			}
		}
		return nil
	})
	titleBoxState.SetExitAction(func() error {
		if keyb.IsKeySPressed(Keyboard.Enter) {
			TitleBox.DeleteLastCharacter() //delete last /n
		}
		TitleBox.OnLeave()
		core.SetVisibilityCursor(false)
		return nil
	})

	bottonSendEditState := StateMachine.CreateBuilderStateBase("BottonsEditState")
	bottonSendEditState.SetEntryAction(func() error {
		SendButton.GetVisibleArea().SetBorderColor(Color.Get(Color.White, Color.None))
		return nil
	})
	bottonSendEditState.SetActionDo(func() error {
		if keyb.IsKeySPressed(Keyboard.Enter) {
			SendButton.OnClick()
		}
		return nil
	})
	bottonSendEditState.SetExitAction(func() error {
		SendButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Gray, Color.None))
		return nil
	})

	bottonCancelEditState := StateMachine.CreateBuilderStateBase("BottonsEditState")
	bottonCancelEditState.SetEntryAction(func() error {
		CancelButton.GetVisibleArea().SetBorderColor(Color.Get(Color.White, Color.None))
		return nil
	})
	bottonCancelEditState.SetActionDo(func() error {
		if keyb.IsKeySPressed(Keyboard.Enter) {
			CancelButton.OnClick()
		}
		return nil
	})
	bottonCancelEditState.SetExitAction(func() error {
		CancelButton.GetVisibleArea().SetBorderColor(Color.Get(Color.Gray, Color.None))
		return nil
	})
	isOk, err := StateMachine.Do(
		//EDITPART
		editState.AddState(titleBoxState),
		editState.AddState(textBoxState),
		editState.AddState(bottonSendEditState),
		editState.AddState(bottonCancelEditState),
		//TODOSTATE
		todoState.AddState(caroselloState),
		todoState.AddState(bottonsCaroselloState),
		//TEXTBOXSTATE
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editState),
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeyPressed('k') && !TextBox.IsTyping() && !TitleBox.IsTyping()
		}, titleBoxState),
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlK)
		}, titleBoxState),
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlJ)
		}, bottonSendEditState),
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeyPressed('k') && !TextBox.IsTyping() && !TitleBox.IsTyping()
		}, titleBoxState),
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlH)
		}, todoState),
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeyPressed('h') && !TextBox.IsTyping() && !TitleBox.IsTyping()
		}, todoState),
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlH)
		}, todoState),
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeyPressed('j') && !TextBox.IsTyping() && !TitleBox.IsTyping()
		}, bottonSendEditState),
		//TITLEBOXSTATE
		textBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlH)
		}, todoState),
		titleBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlJ)
		}, textBoxState),
		titleBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Enter) && TitleBox.IsTyping()
		}, textBoxState),
		titleBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editState),
		titleBoxState.AddBranch(func() bool {
			return keyb.IsKeyPressed('h') && !TitleBox.IsTyping() && !TitleBox.IsTyping()
		}, todoState),
		titleBoxState.AddBranch(func() bool {
			return keyb.IsKeyPressed('j') && !TitleBox.IsTyping() && !TitleBox.IsTyping()
		}, textBoxState),
		titleBoxState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlH)
		}, caroselloState),
		//BOTTONSCAROSSELLOSTATE
		bottonsCaroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlL)
		}, titleBoxState),
		bottonsCaroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc) || keyb.IsKeyPressed('k') || keyb.IsKeyPressed('j')
		}, caroselloState),
		bottonsCaroselloState.AddBranch(func() bool {
			return carosello.GetElementsNumber() == 0
		}, editState),
		//BOTTONSENDEDITSTATE
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeyPressed('l')
		}, bottonCancelEditState),
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editState),
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeyPressed('h')
		}, todoState),
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlH)
		}, caroselloState),
		bottonSendEditState.AddBranch(func() bool {
			return keyb.IsKeyPressed('k')
		}, textBoxState),
		//BOTTONCANCELEDITSTATE
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeyPressed('h')
		}, bottonSendEditState),
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlH)
		}, caroselloState),
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Esc)
		}, editState),
		bottonCancelEditState.AddBranch(func() bool {
			return keyb.IsKeyPressed('k')
		}, textBoxState),
		//CAROSSELLOSTATE
		caroselloState.AddBranch(func() bool {
			return carosello.GetElementsNumber() == 0
		}, editState),
		caroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlL)
		}, titleBoxState),
		caroselloState.AddBranch(func() bool {
			return keyb.IsKeyPressed('l')
		}, editState),
		caroselloState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.Enter) && carosello.GetElementsNumber() > 0 && !isInDoneState
		}, bottonsCaroselloState),
		//STATEShow
		showTodoState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlD)
		}, showDoneState),
		showDoneState.AddBranch(func() bool {
			return keyb.IsKeySPressed(Keyboard.CtrlD)
		}, showTodoState),

		stateMachine.AddBuilderComposite(editState),
		stateMachine.AddBuilderComposite(todoState),
		stateMachine.AddBuilder(titleBoxState),
		stateMachine.AddBuilder(showTodoState),

		stateMachine.Start(),
		EventManager.Subscribe(ClockEvent, 500, func(_ []any) {
			keyb.CleanKeyboardState()
			stateMachine.Clock()
		}),
	)

	if !isOk {
		panic(err)
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
	stateMachine.Clock()
	core.SetCur(x, y)
	return true
}
