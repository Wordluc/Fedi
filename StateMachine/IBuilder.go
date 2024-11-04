package StateMachine

import State "Fedi/StateMachine/internal"

type IBuilder interface {
	Build() (State.IState, error)
	GetInstance() State.IState
}
