package State

import "errors"

type StateEnd struct {
	StateName         string
	TransitionTo      Transition
	HeadsStateMachine *HeadsStateMachine

	IEntryAction func() error
	IExitAction  func() error
	IDoAction    func() error
}

func (s *StateEnd) EntryAction() error {
	if s.IEntryAction == nil {
		return nil
	}
	return s.IEntryAction()
}

func (s *StateEnd) ExitAction() error {
	if s.HeadsStateMachine == nil {
		return errors.New("no state machine")
	}
	if s.IExitAction == nil {
		return nil
	}
	return s.IExitAction()
}

func (s *StateEnd) DoAction() error {
	if s.IDoAction == nil {
		return nil
	}
	return s.IDoAction()
}

func (s *StateEnd) SetHeadsStateMachine(headsStateMachine *HeadsStateMachine) {
	s.HeadsStateMachine = headsStateMachine
}

func (s *StateEnd) CheckTransition() error {
	ok, err := s.TransitionTo.TryTransition()
	if err != nil {
		return err
	}
	if ok {
		s.HeadsStateMachine.RemoveHead(s)
	}
	return nil
}
