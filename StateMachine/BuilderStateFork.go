package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"
)

type tupleBuilderCond struct {
	builder IBuilder
	cond    func() bool
}
type BuilderStateFork struct {
	state *State.StateFork
	builders []IBuilder
	cond func () bool
	alreadyBuilt bool
}

func CreateBuilderStateFork(nameState string) *BuilderStateFork {
	return &BuilderStateFork{
		state: &State.StateFork{
			StateName: nameState,
		},
	}
}

func (b *BuilderStateFork) SetEntryAction(entryAction func() error) *BuilderStateFork {
	b.state.IEntryAction = entryAction
	return b
}

func (b *BuilderStateFork) SetExitAction(exitAction func() error) *BuilderStateFork {
	b.state.IExitAction = exitAction
	return b
}

func (b *BuilderStateFork) SetActionDo(do func() error) *BuilderStateFork {
	b.state.IDoAction = do
	return b
}

func (b *BuilderStateFork) GetInstance() State.IState {
	return b.state
}

func (b *BuilderStateFork) Build() (State.IState, error) {
	if b.alreadyBuilt {
		return b.state, nil
	}
	if b.state == nil {
		return nil, errors.New("no state")
	}
	if b.state.StateName == "" {
		return nil, errors.New("no state name")
	}
	for _, t := range b.builders {
		if to, e := t.Build(); e != nil {
			return nil, e
		} else {
			b.state.Transitions = append(b.state.Transitions, State.CreateTransition(b.state, to, b.cond))
		}
	}
	if n := len(b.state.Transitions); n == 0 {
		return nil, errors.New("no transition")
	}
	for _, transition := range b.state.Transitions {
		if e := transition.IsValid(); e != nil {
			return nil, e
		}
	}
	b.alreadyBuilt = true
	return b.state, nil
}

func (b *BuilderStateFork) AddTos(cond func() bool, tos ... IBuilder)error {
	for _, to := range tos {
		if _,ok:=to.(*BuilderStateMerge);ok{
			return errors.New("cannot use merge as next")
		}
	}
	b.cond = cond
	b.builders = tos
	return nil
}
