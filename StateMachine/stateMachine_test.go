package StateMachine

import (
	"testing"
)
//a->b->c
func TestStateMachineLinearState(t *testing.T) {
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
   end,e:=buildEnd.Build()
	if e!=nil{
		t.Fatalf(e.Error())
	}
	builderB.SetTransition(func() bool {
		return true
	},end)

	b,e:=builderB.Build()
	if e!=nil{
		t.Fatalf(e.Error())
	}
	builderA.SetTransition(func() bool {
		return true
	},b)
	a,e:=builderA.Build()
	if e!=nil{
		t.Fatalf(e.Error())
	}
	stateMachine.AddHead(a)
	stateMachine.Clock()
	if res!="entryA->doA->exitA->entryB->"{
		t.Fail()
	}
	stateMachine.Clock()
	if res!="entryA->doA->exitA->entryB->doB->exitB->"{
		t.Fail()
	}
}

