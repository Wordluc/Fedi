package State

type StateBase struct {
	StateName         string
	TransitionTo      []*Transition
	HeadsStateMachine *HeadsStateMachine

	IEntryAction func() error
	IExitAction  func() error
	IDoAction    func() error
	nCallTimes    int
}

func (s *StateBase) EntryAction() error {
	if s.IEntryAction == nil {
		return nil
	}
	s.nCallTimes=0
	return s.IEntryAction()
}
func (s *StateBase) GetTransitionsTo() []*Transition {
	return s.TransitionTo
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

func (s *StateBase) CheckTransition() (bool,error) {
	for _, transition := range s.TransitionTo {
		ok, err := transition.TryTransition()
		if err != nil {
			return false,err
		}
		if ok {
			s.HeadsStateMachine.RemoveHead(s)
			s.HeadsStateMachine.AddHead(transition.to)
			return true,nil
		}
	}
	if len(s.TransitionTo) == 0 {
		s.HeadsStateMachine.RemoveHead(s)
		return true,nil
	}
	return false,nil
}
