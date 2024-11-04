package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"
)

type BuilderStateEnd struct {
	state *State.StateEnd
}

func CreateBuilderStateEnd(nameState string) *BuilderStateEnd {
	return &BuilderStateEnd{
		state: &State.StateEnd{
			StateName: nameState,
		},
	}
}

func (b *BuilderStateEnd) SetEntryAction(entryAction func() error) *BuilderStateEnd {
	b.state.IEntryAction = entryAction
	return b
}

func (b *BuilderStateEnd) SetExitAction(exitAction func() error) *BuilderStateEnd {
	b.state.IExitAction = exitAction
	return b
}

func (b *BuilderStateEnd) SetActionDo(do func() error) *BuilderStateEnd {
	b.state.IDoAction = do
	return b
}

func (b *BuilderStateEnd) Build() (State.IState,error) {
	if b.state==nil{
		return nil,errors.New("no state")
	}
	if b.state.StateName==""{
		return nil,errors.New("no state name")
	}
	return b.state,nil
}
