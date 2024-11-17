package StateMachine

import (
	State "Fedi/StateMachine/internal"
)

type StateMachine struct {
	stateToAdd State.IState
	heads *State.HeadsStateMachine
}

func CreateStateMachine() *StateMachine {
	return &StateMachine{
		heads: &State.HeadsStateMachine{Heads: make([]State.IState, 0)},
	}
}

func (m *StateMachine) Clock()error {
	for _, head := range m.heads.GetHeads() {
		head.SetHeadsStateMachine(m.heads)
		pass,e:=head.CheckTransition()
		if e!=nil{
			return e
		}
		if pass{
			continue
		}
		e=head.DoAction()
		if e!=nil{
			return e
		}
	}
	return nil
}

func (m *StateMachine) AddBuilder(state IBuilder)error {
	build, err := state.Build()
	if err != nil {
		return err
	}
	m.heads.AddHead(build)
	return nil
}
