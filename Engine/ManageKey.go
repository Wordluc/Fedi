package Engine

import (
	"TUI/Devices/KeyBoard"
	ITerminal "TUI/Devices/Terminal"
	"errors"
	"slices"
)

type Core struct {
	Keyb   KeyBoard.IKeyBoard
	Term   ITerminal.ITerminal
	events []Event
}

func Setup(pkey KeyBoard.IKeyBoard, pterm ITerminal.ITerminal) (*Core, error) {
	core := &Core{}
	if pkey == nil || pterm == nil {
		return core, errors.New("KeyBoard or Terminal is nil")
	}
	core.Keyb = pkey
	core.Term = pterm
	core.events = make([]Event, 0)
	return core, nil
}

func (core *Core) AddEvent(event Event) error {
	if slices.Contains(core.events, event) {
		return errors.New("Faild to add event, event already exists")
	}
	(*core).events = append(core.events, event)
	return nil
}

func (core *Core) LoopEvent() {
	for _, event := range core.events {
		if event == nil {
			break
		}
		event.EventHandler(core)
	}
}
