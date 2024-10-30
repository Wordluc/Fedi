package StateMachine

import "errors"

type HeadsStateMachine struct {
	heads []IState
}

func (h *HeadsStateMachine) AddHead(state IState) {
	h.heads=append(h.heads,state)
}

func (h *HeadsStateMachine) RemoveHead(state IState) {
	for i := 0; i < len(h.heads); i++ {
		if h.heads[i] == state {
			h.heads = append(h.heads[:i], h.heads[i+1:]...)
		}
	}
}

func (h *HeadsStateMachine) GetHeads() []IState {
	return h.heads
}

type IState interface {
	EntryAction() error
	ExitAction() error
	DoAction() error
	CheckTransition() (bool, error)
}

type StateBase struct {
	stateName         string
	transitionTo      Transition
	awaiting          bool
	done              bool
	headsStateMachine *HeadsStateMachine

	entryAction func() error
	exitAction  func() error
	doAction    func() error
}

func (s *StateBase) EntryAction() error {
	if s.entryAction == nil || s.entryAction() == nil {
		s.awaiting = false
		s.headsStateMachine.AddHead(s)
	}
	return nil
}

func (s *StateBase) ExitAction() error {
	if s.headsStateMachine == nil {
		return errors.New("no state machine")
	}
	if s.exitAction == nil || s.exitAction() == nil {
		s.done = true
		s.headsStateMachine.RemoveHead(s)
	}
	return nil
}

func (s *StateBase) DoAction() error {
	if s.doAction == nil {
		return nil
	}
	return s.doAction()
}

func (s *StateBase) CheckTransition() (bool, error) {
	return s.transitionTo.TryTransition()
}
