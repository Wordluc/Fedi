package State

import "errors"

type StateComposite struct {
	NameState string
	InternalStates []IState
	TransitionsTo    []*Transition
	HeadsStateMachine *HeadsStateMachine
	IEntryAction func() error
	IExitAction  func() error
	IDoAction    func() error
}
func(s *StateComposite)	EntryAction() error{
	if s.IEntryAction == nil {
		return nil
	}
	return s.IEntryAction()
}
func(s *StateComposite)	ExitAction() error{
	for _, state := range s.InternalStates {
		if s.HeadsStateMachine==nil{
			break
		}
		if s.HeadsStateMachine.isActive(state) {
			state.ExitAction()
		}
	}
	if s.IExitAction != nil {
		s.IExitAction()
	}
	return nil
}
func(s *StateComposite)	DoAction() error{
	if s.IDoAction == nil {
		return nil
	}
	return s.IDoAction()
}
func(s *StateComposite)	SetHeadsStateMachine(headsStateMachine *HeadsStateMachine){
	s.HeadsStateMachine = headsStateMachine
}
func(s *StateComposite)	GetTransitionsTo() []*Transition{
	return s.TransitionsTo
}
func (s *StateComposite) AddInternalState(state IState) {
	s.InternalStates = append(s.InternalStates, state)
}
func (s *StateComposite) GetHead()(IState,error) {
	if len(s.InternalStates)==0{
		return nil,errors.New("no state")
	}
	return s.InternalStates[0],nil
}
func(s *StateComposite)	CheckTransition() (bool,error){
	panic("should not be called")
}

func (s *StateComposite) isInside(state IState) bool {
	for _, internalState := range s.InternalStates {
		if internalState == state {
			return true
		}
	}
	return false
}
