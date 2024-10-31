package State
type HeadsStateMachine struct {
	Heads []IState
}

func (h *HeadsStateMachine) AddHead(state IState) {
	h.Heads = append(h.Heads, state)
	state.EntryAction()
}

func (h *HeadsStateMachine) RemoveHead(state IState) {
	for i := 0; i < len(h.Heads); i++ {
		if h.Heads[i] == state {
			state.ExitAction()
			h.Heads = append(h.Heads[:i], h.Heads[i+1:]...)
		}
	}
}

func (h *HeadsStateMachine) GetHeads() []IState {
	return h.Heads
}
