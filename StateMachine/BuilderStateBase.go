package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"

)

type BuilderStateBase struct {
	state *State.StateBase
	tos []*tupleBuilderCond
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

	for _, t := range b.tos {
		if t.builder == nil {
			return nil, errors.New("no builder")
		}
		if t.cond == nil {
			return nil, errors.New("no condition")
		}
		if to, e := t.builder.Build(); e != nil {
			return nil, e
		} else {
			b.state.TransitionTo = append(b.state.TransitionTo, State.CreateTransition(b.state, to, t.cond))
		}
	}

	return b.state,nil
}

func (b *BuilderStateBase) AddBranch(cond func() bool,builderNext IBuilder)error {
	if _,ok:=builderNext.(*BuilderStateMerge);ok{
		return errors.New("invalid next state, cannot use merge as next, use wait instead")
	}
	b.tos = append(b.tos, &tupleBuilderCond{
		cond:    cond,
		builder: builderNext,
	})
	return nil
}
