package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"
)

type BuilderStateMerge struct {
	state *State.StateMerge
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

func (b *BuilderStateMerge) AddInTransition(cond func() bool,from State.IState) error {
	state,ok:=from.(*State.StateBase)
	if !ok{
		return errors.New("expected StateBase")
	}
	if state.TransitionTo.IsValid()==nil{
		return errors.New("transition already set")
	}
	state.TransitionTo = *State.CreateTransition(state, b.state, cond)
	return nil
}

func (b *BuilderStateMerge) Build() (*State.StateMerge,error) {
	if b.state==nil{
		return nil,errors.New("no state")
	}
	if b.state.StateName==""{
		return nil,errors.New("no state name")
	}
	if e:= b.state.TransitionTo.IsValid();e!=nil{
		return nil,e
	}
	if n:= len(b.state.InTransitions);n==0{
		return nil,errors.New("no inTransitions")
	}
	return b.state,nil
}

