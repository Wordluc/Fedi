package State

import (
	"errors"
)

type StateMerge struct {
	StateName         string
	InTransitions     []Transition
	TransitionTo      Transition
	HeadsStateMachine *HeadsStateMachine

	IEntryAction func() error
	IExitAction  func() error
	IDoAction    func() error
}

func (s *StateMerge) EntryAction() error {
	if s.IEntryAction == nil {
		return s.IEntryAction()
	}
	return nil
}

func (s *StateMerge) ExitAction() error {
	if s.HeadsStateMachine == nil {
		return errors.New("no state machine")
	}
	if s.IExitAction == nil || s.IExitAction() == nil {
		s.HeadsStateMachine.RemoveHead(s)
	}
	return nil
}

func (s *StateMerge) DoAction() error {
	if s.IDoAction == nil {
		return nil
	}
	return s.DoAction()
}

func (s *StateMerge) SetHeadsStateMachine(headsStateMachine *HeadsStateMachine) {
	s.HeadsStateMachine = headsStateMachine
}

func (s *StateMerge) CheckTransition() error {
	ok, err := s.TransitionTo.TryTransition()
	if err != nil {
		return err
	}
	if ok {
		s.HeadsStateMachine.RemoveHead(s)
		s.HeadsStateMachine.AddHead(s.TransitionTo.to)
	}
	return err
}
