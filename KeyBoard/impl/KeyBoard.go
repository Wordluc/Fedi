package impl

import (
	"TUI/Terminal/Token"
	"github.com/eiannone/keyboard"
)

type loop func() bool
type KeyBoard struct {
	key  stateKey
	loop loop
}
type stateKey struct {
	key  keyboard.Key
	rune rune
}

func (t *KeyBoard) Start(loop loop) error {
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
		print("dio cane")
	}
	return nil
}

func (t *KeyBoard) Stop() {
	keyboard.Close()
}

func (t *KeyBoard) GetKey() (byte, error) {
	return byte(t.key.rune), nil
}

func (t *KeyBoard) IsPressed(token Token.Token) bool {
	key := t.key.key
	if v, e := mapTokenToKey(token); e == nil {
		if v == key {
			return true
		}
	}
	return false
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
