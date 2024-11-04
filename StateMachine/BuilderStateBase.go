package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"

)

type BuilderStateBase struct {
	state *State.StateBase
	builderNext IBuilder
	conditionNext func() bool
}

func CreateBuilderStateBase(nameState string) *BuilderStateBase {
	return &BuilderStateBase{
		state: &State.StateBase{
			StateName: nameState,
		},
	}
}
func (b *BuilderStateBase) GetInstance() State.IState {
	return b.state
}
func (b *BuilderStateBase) SetEntryAction(entryAction func() error) *BuilderStateBase {
	b.state.IEntryAction = entryAction
	return b
}

func (b *BuilderStateBase) SetExitAction(exitAction func() error) *BuilderStateBase {
	b.state.IExitAction = exitAction
	return b
}

func (b *BuilderStateBase) SetActionDo(do func() error) *BuilderStateBase {
	b.state.IDoAction = do
	return b
}

func (b *BuilderStateBase) Build() (State.IState,error) {
	if b.state==nil{
		return nil,errors.New("no state")
	}

	if b.state.StateName==""{
		return nil,errors.New("no state name")
	}

	if b.builderNext==nil{
		return nil,nil
	}

	to, e:= b.builderNext.Build()
	if e!=nil{
		return nil,e
	}
	b.state.TransitionTo = *State.CreateTransition(b.state, to, b.conditionNext)

	if e:= b.state.TransitionTo.IsValid();e!=nil{
		return nil,e
	}
	return b.state,nil
}

func (b *BuilderStateBase) SetNext(cond func() bool,builderNext IBuilder) {
	b.builderNext = builderNext
	b.conditionNext = cond
}
