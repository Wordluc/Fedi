package State

import "errors"

type Transition struct {
	from      IState
	to        IState
	condition func() bool
}

func (t *Transition) TryTransition() (bool, error) {
	if t.condition == nil {
		return false, errors.New("no condition")
	}
	return t.condition(), nil
}
func (t *Transition) IsValid()error {
	
	if t.from==nil{
		return errors.New("no from state")
	}
	if t.to==nil{
		return errors.New("no to state")
	}
	if t.condition==nil{
		return errors.New("no condition")
	}
	return nil
}

func CreateTransition(from, to IState, condition func() bool) *Transition {
	return &Transition{
		to:        to,
		from:      from,
		condition: condition,
	}
}
