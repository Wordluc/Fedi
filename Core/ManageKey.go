package Core

import (
	"TUI/KeyBoard"
	"TUI/Terminal"
	"errors"
	"slices"
)

var keyb *KeyBoard.IKeyBoard = nil
var term *Terminal.Terminal
var events []Event

type Event interface {

}

type ArrowEvent struct {

}

func Setup(pkey *KeyBoard.IKeyBoard, pterm *Terminal.Terminal) error {
	if pkey == nil || pterm == nil {
		return errors.New("Faild to setup")
	}
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
