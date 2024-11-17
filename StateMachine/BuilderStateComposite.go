package StateMachine

import (
	State "Fedi/StateMachine/internal"
	"errors"
)
type StateCompositeBuilder struct {
	builder *State.StateComposite
	buildersInternalState []IBuilder
	tos []*tupleBuilderCond
	alreadyBuilt bool
}
func CreateBuilderStateComposite(nameState string) *StateCompositeBuilder {
	return &StateCompositeBuilder{
		builder: &State.StateComposite{
			NameState: nameState,
			InternalStates: make([]State.IState, 0),
		},
	}
}
func (b *StateCompositeBuilder) Build() (State.IState, error) {
	if b.alreadyBuilt {
		return b.builder, nil
	}
	for _, t := range b.buildersInternalState {
		if to, e := t.Build(); e != nil {
			return nil, e
		} else {
			b.builder.InternalStates = append(b.builder.InternalStates, to)
		}
	}
	b.alreadyBuilt = true

	for _, t := range b.tos {
		if t.builder == nil {
			return nil, errors.New("no builder")
		}
		if t.cond == nil {
			return nil, errors.New("no condition")
		}
		if to, e := t.builder.Build(); e != nil {
			return nil, e
		} else {
			b.builder.TransitionsTo = append(b.builder.TransitionsTo, State.CreateTransition(b.builder, to, t.cond))
		}
	}
	return b.builder, nil
}
func (b *StateCompositeBuilder) GetInstance() State.IState {
	return b.builder
}

func (b *StateCompositeBuilder) AddState(state IBuilder) error {
	b.buildersInternalState = append(b.buildersInternalState, state)
	return nil
}
func (b *StateCompositeBuilder) SetEntryAction(entryAction func() error) *StateCompositeBuilder {
	b.builder.IEntryAction = entryAction
	return b
}
func (b *StateCompositeBuilder) SetExitAction(exitAction func() error) *StateCompositeBuilder {
	b.builder.IExitAction = exitAction
	return b
}
func (b *StateCompositeBuilder) SetActionDo(do func() error) *StateCompositeBuilder {
	b.builder.IDoAction = do
	return b
}
