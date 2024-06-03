package Token

type Token int
type Arrow Token
const (
	Arrow_Left  Arrow = 0
	Arrow_Right       = 1
	Arrow_Up          = 2
	Arrow_Down        = 3
)

func (t Token) String() string {
	return [...]string{
		"Arrow_Left",
		"Arrow_Right",
		"Arrow_Up",
		"Arrow_Down",
	}[t]
}
