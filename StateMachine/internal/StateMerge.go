package State

import (
	"errors"
)

type StateMerge struct {
	StateName         string
	ToWait     			[]IState
	TransitionTo      Transition
	HeadsStateMachine *HeadsStateMachine

	IEntryAction func() error
	IExitAction  func() error
	IDoAction    func() error
}

func (s *StateMerge) EntryAction() error {
	if s.IEntryAction == nil {
		return nil
	}
	return s.IEntryAction()
}

func (s *StateMerge) ExitAction() error {
	if s.HeadsStateMachine == nil {
		return errors.New("no state machine")
	}
	if s.IExitAction == nil {
		return nil
	}
	return s.IExitAction()
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

func (s *StateMerge) GetTransitionsTo() []Transition {
	return []Transition{s.TransitionTo}
}
func (s *StateMerge) CheckTransition() error {
	for _, state := range s.ToWait {
		tran:=state.GetTransitionsTo()
		for _, transition := range tran {
			if transition.to!=s{
				return errors.New("invalid transition: cannot transition to a state that doesn't merge back")
			}
			if !transition.IsDone() {
				return nil
			}
		}
	}
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
