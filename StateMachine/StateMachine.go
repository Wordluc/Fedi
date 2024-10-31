package StateMachine

import State "Fedi/StateMachine/internal"

type StateMachine struct {
	begin State.IState
	heads *State.HeadsStateMachine
}

func CreateStateMachine() *StateMachine {
	return &StateMachine{
		heads: &State.HeadsStateMachine{Heads: make([]State.IState, 0)},
	}
}

func (m *StateMachine) Clock() {
	for _, head := range m.heads.GetHeads() {
		head.SetHeadsStateMachine(m.heads)
		head.DoAction()
		head.CheckTransition()
	}
}

func (m *StateMachine) AddHead(state State.IState) {
	m.heads.AddHead(state)
}
