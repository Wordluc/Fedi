package impl

import (
	"TUI/Devices/KeyBoard"
	"TUI/Devices/Terminal/Token"

	"github.com/eiannone/keyboard"
)

type ImplKeyBoard struct {
	key  stateKey
	loop KeyBoard.Loop
}
type stateKey struct {
	key  keyboard.Key
	rune rune
}

func (t *ImplKeyBoard) Start(loop KeyBoard.Loop) error {
	keyboard.Open()
	eventKey, e := keyboard.GetKeys(10)
	if e != nil {
		return e
	}
	t.loop = loop
	for {
		v := <-eventKey
		t.key = stateKey{key: v.Key, rune: v.Rune}
		if t.loop() {
			break
		}
	}
	return nil
}

func (t *ImplKeyBoard) Stop() {
	keyboard.Close()
}

func (t *ImplKeyBoard) GetKey() (byte, error) {
	return byte(t.key.rune), nil
}

func (t *ImplKeyBoard) IsTokenPressed(token Token.Token) bool {
	key := t.key.key
	if v, e := mapTokenToKey(token); e == nil {
		if v == key {
			return true
		}
	}
	return false
}

func (t *ImplKeyBoard) IsKeyPressed(key byte) bool {
	return byte(t.key.rune)==key 
}

func mapTokenToKey(token Token.Token) (keyboard.Key, error) {
	switch token {
	case Token.Arrow_Left:
		return keyboard.KeyArrowLeft, nil
	case Token.Arrow_Right:
		return keyboard.KeyArrowRight, nil
	case Token.Arrow_Up:
		return keyboard.KeyArrowUp, nil
	case Token.Arrow_Down:
		return keyboard.KeyArrowDown, nil
	}
	return 0, nil
}
