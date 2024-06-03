package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"golang.org/x/term"
	"os"
)

type Terminal struct {
	term *term.State
}

func (t *Terminal) Start() error {
	tt, e := term.MakeRaw(int(os.Stdin.Fd()))
	if e != nil {
		return e
	}
	t.term = tt
	term.NewTerminal(os.Stdin, "")
	return nil
}

func (t *Terminal) Stop() {
	term.Restore(int(os.Stdin.Fd()), t.term)
}

func main() {
	t := Terminal{}
	t.Start()
	keyboard.Open()
	defer t.Stop()
	os.Stdout.Write([]byte("\033[H\033[2J"))
	for {

		var a, b, _ = keyboard.GetKey()
		fmt.Println(a, b)
		if b == keyboard.KeyCtrlW {
			break
		}
		os.Stdout.Write([]byte(string(a))) // write a )
	}
}
