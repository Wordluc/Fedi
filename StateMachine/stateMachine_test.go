package StateMachine

import (
	"testing"
)
//a->b->c
func TestStateMachineLinearStates(t *testing.T) {
	stateMachine := CreateStateMachine()
	var res string =""
	builderA:=CreateBuilderStateBase("a")
	builderA.SetActionDo(func() error {
		res+="doA->"
		return nil
	})
	builderA.SetExitAction(func() error {
		res+="exitA->"
		return nil
	})
	builderA.SetEntryAction(func() error {
		res+="entryA->"
		return nil
	})
	builderB:=CreateBuilderStateBase("b")
	builderB.SetActionDo(func() error {
		res+="doB->"
		return nil
	})
	builderB.SetExitAction(func() error {
		res+="exitB->"
		return nil
	})
	builderB.SetEntryAction(func() error {
		res+="entryB->"
		return nil
	})

	buildEnd:=CreateBuilderStateEnd("end")
	builderA.SetNext(func() bool {return true },builderB)
	builderB.SetNext(func() bool {return true },buildEnd)
	stateMachine.AddBuilder(builderA)
	stateMachine.Clock()
	if res!="entryA->doA->exitA->entryB->"{
		t.Fail()
	}
	stateMachine.Clock()
	if res!="entryA->doA->exitA->entryB->doB->exitB->"{
		t.Fail()
	}
}

