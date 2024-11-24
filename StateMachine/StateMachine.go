package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"
)

type StateMachine struct {
	stateToAdd State.IState
	heads *State.HeadsStateMachine
	toActivate []State.IState
}

func CreateStateMachine() *StateMachine {
	return &StateMachine{
		heads: &State.HeadsStateMachine{State: make([]State.IState, 0)},
	}
}
func (m *StateMachine) Start()error {
	if m.toActivate==nil{
		return nil
	}
	for _, state := range m.toActivate {
		m.heads.AddHead(nil,state)
	}
	m.toActivate = nil
	return nil
}

func (m *StateMachine) Clock()error {
	if m.toActivate!=nil{
		return errors.New("Machine is not started")
	}
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
	m.toActivate = append(m.toActivate, build)
	return nil
}

func (m *StateMachine) AddBuilderComposite(state *StateCompositeBuilder)error {
	build, err := state.Build()
	if err != nil {
		return err
	}
	m.heads.WaitingStatesComposite = append(m.heads.WaitingStatesComposite, build.(*State.StateComposite))
	return nil
}
