package Core

import (
	"TUI/KeyBoard"
	"TUI/Terminal"
	"TUI/Terminal/Token"
	"errors"
	"slices"
)

var keyb KeyBoard.IKeyBoard = nil
var term Terminal.ITerminal
var events []Event

type Event interface {
    EventHandler()
}
type Handler func(Key Token.Token)

type ArrowEvent struct {
   Key []Token.Token
	 Handler Handler
}

func (event ArrowEvent) EventHandler() {
	for _, key := range event.Key {
		if(keyb.IsPressed(key)){
			event.Handler(key)
		}
	}
}
func Setup(pkey KeyBoard.IKeyBoard, pterm Terminal.ITerminal) error {
	keyb = pkey
	term = pterm
	return nil
}

func AddEvent(event Event)error {
	if slices.Contains(events, event) {
		return errors.New("Faild to add event, event already exists")
	}
	events = append(events, event)
	return nil
}

func LoopEvent(){
	for _, event := range events {
		event.EventHandler()
	}
}
