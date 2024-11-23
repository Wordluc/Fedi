package State

import "slices"

type HeadsStateMachine struct {
	State                  []IState
	ActiveStatesComposite  []*StateComposite
	WaitingStatesComposite []*StateComposite
}

func (h *HeadsStateMachine) AddHead(state IState) {
	defer func() {
		h.State = append(h.State, state)
		state.EntryAction()
	}()
	if stateC, ok := state.(*StateComposite); ok {
		var err error
		state, err = stateC.GetHead()
		if err != nil {
			return
		}
	}
	for i, stateComposite := range h.WaitingStatesComposite {
		if stateComposite==nil{
			continue
		}
		if stateComposite.isInside(state) {
			stateComposite.EntryAction()
			h.ActiveStatesComposite = append(h.ActiveStatesComposite, stateComposite)
			h.WaitingStatesComposite = slices.Delete(h.WaitingStatesComposite, i, i+1)
		}
	}
	for i, stateComposite := range h.ActiveStatesComposite {
		if stateComposite==nil{
			continue
		}
		if !stateComposite.isInside(state) {
			stateComposite.ExitAction()
			h.WaitingStatesComposite = append(h.WaitingStatesComposite, stateComposite)
			h.ActiveStatesComposite = slices.Delete(h.ActiveStatesComposite, i, i+1)
		}
	}
}

func (h *HeadsStateMachine) RemoveHead(state IState) {
	for i := 0; i < len(h.State); i++ {
		if h.State[i] == state {
			state.ExitAction()
			h.State = slices.Delete(h.State, i, i+1)
			if comp, ok := state.(*StateComposite); ok {
				h.WaitingStatesComposite = append(h.WaitingStatesComposite, comp)
				h.ActiveStatesComposite = slices.Delete(h.ActiveStatesComposite, i, i+1)
			}
		}
	}
}

func (h *HeadsStateMachine) GetHeads() []IState {
	return h.State
}
