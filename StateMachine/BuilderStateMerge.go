package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"
)

type BuilderStateMerge struct {
	state   *State.StateMerge
	toWaits []IBuilder
	to      tupleBuilderCond
}

func CreateBuilderStateMerge(nameState string) *BuilderStateMerge {
	return &BuilderStateMerge{
		state: &State.StateMerge{
			StateName: nameState,
		},
	}
}

func (b *BuilderStateMerge) SetEntryAction(entryAction func() error) *BuilderStateMerge {
	b.state.IEntryAction = entryAction
	return b
}

func (b *BuilderStateMerge) SetExitAction(exitAction func() error) *BuilderStateMerge {
	b.state.IExitAction = exitAction
	return b
}

func (b *BuilderStateMerge) SetActionDo(do func() error) *BuilderStateMerge {
	b.state.IDoAction = do
	return b
}

func (b *BuilderStateMerge) GetInstance() State.IState {
	return b.state
}

func (b *BuilderStateMerge) Build() (State.IState, error) {
	if b.state == nil {
		return nil, errors.New("no state")
	}
	if b.state.StateName == "" {
		return nil, errors.New("no state name")
	}
	for _, t := range b.toWaits {
		istance := t.GetInstance()
		b.state.ToWait = append(b.state.ToWait,istance)
	}

	if b.to.builder == nil {
		return nil, errors.New("no builder")
	}
	if b.to.cond == nil {
		return nil, errors.New("no condition")
	}
	if to, e := b.to.builder.Build(); e != nil {
		return nil, e
	} else {
		b.state.TransitionTo = *State.CreateTransition(b.state, to, b.to.cond)
	}
	if e := b.state.TransitionTo.IsValid(); e != nil {
		return nil, e
	}
	return b.state, nil
}

func (b *BuilderStateMerge) SetNext(condToOut func() bool, outBuild IBuilder) {
	b.to = tupleBuilderCond{
		cond:    condToOut,
		builder: outBuild,
	}
}

func (b *BuilderStateMerge) AddToWait(cond func() bool, toWait IBuilder) error {
	if state,ok:=toWait.(*BuilderStateBase);ok{
		state.builderNext = b
		state.conditionNext = cond
		b.toWaits = append(b.toWaits, toWait)
		return nil
	}
	return errors.New("cannot wait for this type of state")
}
