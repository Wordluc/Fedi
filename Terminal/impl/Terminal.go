package impl

import (
	"os"

	"golang.org/x/term"
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

func (t *Terminal) Clear() {
	os.Stdout.Write([]byte("\033[H\033[2J"))
}

func (t *Terminal) Print(byte []byte) {
	os.Stdout.Write(byte)
}

func (t *Terminal) PrintStr(str string) {
	os.Stdout.Write([]byte(str))
}
