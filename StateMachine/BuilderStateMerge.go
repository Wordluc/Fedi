package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"
)

type BuilderStateMerge struct {
	state *State.StateMerge
	inBuilds []tupleBuilderCond
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

func (b *BuilderStateMerge) Build() (*State.StateMerge,error) {
	if b.state==nil{
		return nil,errors.New("no state")
	}
	if b.state.StateName==""{
		return nil,errors.New("no state name")
	}
	for _, t := range b.inBuilds {
		if t.builder == nil {
			return nil, errors.New("no builder")
		}
		if t.cond == nil {
			return nil, errors.New("no condition")
		}
		if to, e := t.builder.Build(); e != nil {
			return nil, e
		} else {
			b.state.InTransitions = append(b.state.InTransitions, *State.CreateTransition(to, b.state, t.cond))
		}
	}
	if e:= b.state.TransitionTo.IsValid();e!=nil{
		return nil,e
	}
	if n:= len(b.state.InTransitions);n==0{
		return nil,errors.New("no inTransitions")
	}
	return b.state,nil
}

//add builder state who will merge in this state
//condToIn: condition to SetAddIn
//inBuild: state who will merge
func (b *BuilderStateMerge) SetAddIn(condToIn func () bool,inBuild IBuilder){
	b.inBuilds = append(b.inBuilds,tupleBuilderCond{
		cond:condToIn,
		builder:inBuild,
	})
}
