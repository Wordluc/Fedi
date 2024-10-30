package StateMachine

import "errors"

type Transition struct {
	To        IState
	From      IState
	Condition func() bool
}

func (t *Transition) TryTransition() (bool, error) {
	if t.Condition == nil {
		return false, errors.New("no condition")
	}
	if t.Condition() {
		t.From.ExitAction()
		t.To.EntryAction()
		return true, nil
	}
	return false, nil
}

func CreateTransition(from, to IState, condition func() bool) *Transition {
	return &Transition{
		To:        to,
		From:      from,
		Condition: condition,
	}
}
