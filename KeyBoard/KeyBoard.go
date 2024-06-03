package KeyBoard

import (

	"github.com/eiannone/keyboard"
)

type KeyBoard struct {
	//Handlers map[Token.Token]Handler da spostare in un altro modulo
}

func (t *KeyBoard) Start() {
	keyboard.Open()
}

func (t *KeyBoard) Stop() {
	keyboard.Close()
}

func (t *KeyBoard) GetKey() (byte, error) {
	key, _, err := keyboard.GetKey()
	if err != nil {
		return 0, err
	}
	return byte(key), nil
}

//type ResultHandler interface{}
//type Handler func() (ResultHandler, error)
//
//func (t *KeyBoard) AddHandler(e Handler, tokens ...Token.Token) error {
//	for _,token := range tokens {
//		if _, ok := t.Handlers[token]; ok {
//			return errors.New("Handler for " +  token.String() + " already exists")
//		}
//		t.Handlers[token] = e
//	}
//	return nil
//}
