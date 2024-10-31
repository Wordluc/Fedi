package State

import (
	"errors"
)

type StateFork struct {
	StateName         string
	Transitions       []Transition
	HeadsStateMachine *HeadsStateMachine

	IEntryAction func() error
	IExitAction  func() error
	IDoAction    func() error
}

func (s *StateFork) EntryAction() error {
	if s.IEntryAction == nil {
		return s.IEntryAction() 	
	}
	return nil
}

func (s *StateFork) ExitAction() error {
	if s.HeadsStateMachine == nil {
		return errors.New("no state machine")
	}
	if s.IExitAction == nil || s.IExitAction() == nil {
		s.HeadsStateMachine.RemoveHead(s)
	}
	return nil
}

func (s *StateFork) DoAction() error {
	if s.IDoAction == nil {
		return nil
	}
	return s.DoAction()
}

func (s *StateFork) SetHeadsStateMachine(headsStateMachine *HeadsStateMachine) {
	s.HeadsStateMachine = headsStateMachine
}

func (s *StateFork) CheckTransition() (error) {
	for _, transition := range s.Transitions {
		ok, err := transition.TryTransition()
		if err != nil {
			return err
		}
		if ok {
			s.HeadsStateMachine.RemoveHead(s)
			s.HeadsStateMachine.AddHead(transition.to)
		}
	}
	return nil
}
