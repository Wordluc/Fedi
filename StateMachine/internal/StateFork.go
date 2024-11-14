package State

import (
	"errors"
)

type StateFork struct {
	StateName         string
	Transitions       []*Transition
	HeadsStateMachine *HeadsStateMachine
	passedState int

	IEntryAction func() error
	IExitAction  func() error
	IDoAction    func() error
}

func (s *StateFork) EntryAction() error {
	if s.IEntryAction == nil {
		return nil
	}
	return s.IEntryAction() 	
}

func (s *StateFork) GetTransitionsTo() []*Transition {
	return s.Transitions
}
func (s *StateFork) ExitAction() error {
	if s.HeadsStateMachine == nil {
		return errors.New("no state machine")
	}
	if s.IExitAction == nil {
		return nil
	}
	return s.IExitAction()
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

func (s *StateFork) CheckTransition() (bool,error) {
	for _, transition := range s.Transitions {
		if transition.IsDone() {
			continue
		}
		ok, err := transition.TryTransition()
		if err != nil {
			return false,err
		}
		if ok {
			s.passedState++
			if s.passedState==len(s.Transitions){
				s.HeadsStateMachine.RemoveHead(s)
				return true,nil
			}else{
				s.ExitAction()
			}
			s.HeadsStateMachine.AddHead(transition.to)
		}
	}
	if len(s.Transitions) == 0 {
		s.HeadsStateMachine.RemoveHead(s)
		return true,nil
	}
	return false,nil
}
