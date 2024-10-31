package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"
)

type BuilderStateFork struct {
	state *State.StateFork
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

func (b *BuilderStateFork) Build() (*State.StateFork,error) {
	if b.state==nil{
		return nil,errors.New("no state")
	}
	if b.state.StateName==""{
		return nil,errors.New("no state name")
	}
	if n:= len(b.state.Transitions);n==0{
		return nil,errors.New("no transition")
	}
	for _, transition := range b.state.Transitions {
		if e:= transition.IsValid();e!=nil{
			return nil,e
		}
	}
	return b.state,nil
}

func (b *BuilderStateFork) AddTransition(cond func() bool,to *State.StateBase) {
	b.state.Transitions = append(b.state.Transitions, *State.CreateTransition(b.state, to, cond))
}
