package State

import "errors"

type Transition struct {
	from      IState
	to        IState
	condition func() bool
	isDone    bool
}

func (t *Transition) TryTransition() (bool, error) {
	if t.condition == nil {
		return false, errors.New("no condition")
	}
	t.isDone = t.isDone || t.condition()
	return t.isDone, nil
}

func (t *Transition) IsDone() bool {
	return t.isDone
}

func (t *Transition) GetFrom() IState {
	return t.from
}
func (t *Transition) GetTo() IState {
	return t.to
}
func (t *Transition) IsValid() error {

	if t.from == nil {
		return errors.New("no from state")
	}
	if t.to == nil {
		return errors.New("no 'to' state")
	}
	if t.condition == nil {
		return errors.New("no condition")
	}
	return nil
}

func (t *Transition) SetTo(to IState) {
	t.to = to
}

func CreateTransition(from, to IState, condition func() bool) *Transition {
	return &Transition{
		to:        to,
		from:      from,
		condition: condition,
	}
}
