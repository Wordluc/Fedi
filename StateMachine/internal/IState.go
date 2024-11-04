package State

type IState interface {
	EntryAction() error
	ExitAction() error
	DoAction() error
	SetHeadsStateMachine(headsStateMachine *HeadsStateMachine)
	CheckTransition() (error)
	GetTransitionsTo() []*Transition
}
