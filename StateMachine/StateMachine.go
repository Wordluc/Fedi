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

func (m *StateMachine) AddBuilder(state IBuilder)error {
	build, err := state.Build()
	if err != nil {
		return err
	}
	m.heads.AddHead(build)
	return nil
}
