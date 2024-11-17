package StateMachine

import "strings"
type ErrorComposition struct {
	error string
}
func (e *ErrorComposition) IsEmpty() bool{
	return strings.TrimFunc(e.error, func(r rune) bool {
		return r == '\n' || r == ' ' || r == '\t'
	})==""
}

func (e *ErrorComposition) Add(err error) *ErrorComposition{
	if err != nil {
		e.error =e.error + err.Error() + "\n"
	}
	return e
}
func (e *ErrorComposition) Error() string {
	return e.error+"\n"
}
func Do(possibleErrors ...error) (bool,*ErrorComposition) {
	var error *ErrorComposition = &ErrorComposition{" "}
	for _, e := range possibleErrors {
		if e != nil {
			error.Add(e)
		}
	}
	return error.IsEmpty(),error
}
