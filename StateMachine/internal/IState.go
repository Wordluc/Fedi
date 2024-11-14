package State

type IState interface {
	EntryAction() error
	ExitAction() error
	DoAction() error
	SetHeadsStateMachine(headsStateMachine *HeadsStateMachine)
	CheckTransition() (bool,error)
	GetTransitionsTo() []*Transition
}
