package State

type StateBase struct {
	StateName         string
	TransitionTo      Transition
	HeadsStateMachine *HeadsStateMachine

	IEntryAction func() error
	IExitAction  func() error
	IDoAction    func() error
}

func (s *StateBase) EntryAction() error {
	if s.IEntryAction == nil {
		return nil
	}
	return s.IEntryAction()
}
func (s *StateBase) GetTransitionsTo() []Transition {
	return []Transition{s.TransitionTo}
}
func (s *StateBase) ExitAction() error {
	if s.IExitAction == nil {
		return nil
	}
	return s.IExitAction()
}

func (s *StateBase) DoAction() error {
	if s.IDoAction == nil {
		return nil
	}
	return s.IDoAction()
}

func (s *StateBase) SetHeadsStateMachine(headsStateMachine *HeadsStateMachine) {
	s.HeadsStateMachine = headsStateMachine
}

func (s *StateBase) CheckTransition() (error) {
	ok, err := s.TransitionTo.TryTransition()
	if err != nil {
		return err
	}
	if ok {
		s.HeadsStateMachine.RemoveHead(s)
		s.HeadsStateMachine.AddHead(s.TransitionTo.to)
	}
	return nil
}
