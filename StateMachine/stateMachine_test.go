package StateMachine

import (
	"testing"
)

//a->b->c
func TestStateMachineLinearStates(t *testing.T) {
	stateMachine := CreateStateMachine()
	var trigger=false
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


	builderA.AddBranch(func() bool {return trigger },builderB)
	stateMachine.AddBuilder(builderA)
	stateMachine.Clock()
	if res!="entryA->doA->"{
		t.Fatalf("expected: entryA->doA->, got: %s",res)
	}
	trigger=true
	stateMachine.Clock()
	stateMachine.Clock()
	if res!="entryA->doA->exitA->entryB->exitB->"{
		t.Fatalf("expected: entryA->doA->exitA->entryB->exitB->, got: %s",res)
	}
}

func TestStateMachineFork(t *testing.T) {
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
	builderC:=CreateBuilderStateBase("c")
	builderC.SetActionDo(func() error {
		res+="doC->"
		return nil
	})
	builderC.SetExitAction(func() error {
		res+="exitC->"
		return nil
	})
	builderC.SetEntryAction(func() error {
		res+="entryC->"
		return nil
	})
	
	builderE:=CreateBuilderStateBase("e")
	
	builderE.SetEntryAction(func() error {
		res+="entryE->"
		return nil
	})
	
	builderE.SetExitAction(func() error {
		res+="exitE->"
		return nil
	})
	forkBuilder:=CreateBuilderStateFork("fork")
	
	forkBuilder.SetEntryAction(func() error {
		res+="entryFork->"
		return nil
	})
	forkBuilder.SetExitAction(func() error {
		res+="exitFork->"
		return nil
	})
	mergeBuilder:=CreateBuilderStateMerge("merge")

	builderA.AddBranch(func () bool{return true},forkBuilder)
	forkBuilder.AddTos(func () bool{return true},builderC,builderB)

	error:=mergeBuilder.AddToWait(func () bool{return true},builderC)
	if error!=nil{
		t.Fatalf("build: %s",error)
	}
	error=mergeBuilder.AddToWait(func () bool{return true},builderB)
	if error!=nil{
		t.Fatalf("build: %s",error)
	}

	mergeBuilder.SetNext(func () bool{return true},builderE)
	error=stateMachine.AddBuilder(builderA)
	if error!=nil{
		t.Fatalf("builder: %s",error)
	}
	for i:=0;i<7;i++{
		if e:=stateMachine.Clock();e!=nil{
			t.Fatalf("clock: %s",e)
		}
	}

	exp:="entryA->exitA->entryFork->exitFork->entryC->exitFork->exitC->"
	if res!=exp{
		t.Fatalf("got: %s, expected: %s",res,exp)
	}
}
func TestStateMachineForkWaitWrongState(t *testing.T) {
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
	builderC:=CreateBuilderStateBase("b")
	builderC.SetActionDo(func() error {
		res+="doC->"
		return nil
	})
	builderC.SetExitAction(func() error {
		res+="exitC->"
		return nil
	})
	builderC.SetEntryAction(func() error {
		res+="entryC->"
		return nil
	})
	
	builderE:=CreateBuilderStateBase("e")
	
	builderE.SetEntryAction(func() error {
		res+="entryE->"
		return nil
	})
	
	builderE.SetExitAction(func() error {
		res+="exitE->"
		return nil
	})
	forkBuilder:=CreateBuilderStateFork("fork")
	
	forkBuilder.SetEntryAction(func() error {
		res+="entryFork->"
		return nil
	})
	forkBuilder.SetExitAction(func() error {
		res+="exitFork->"
		return nil
	})
	mergeBuilder:=CreateBuilderStateMerge("merge")

	builderA.AddBranch(func () bool{return true},forkBuilder)
	forkBuilder.AddTos(func () bool{return true},builderC)
	forkBuilder.AddTos(func () bool{return true},builderB)

	mergeBuilder.AddToWait(func () bool{return true},builderA)

	mergeBuilder.SetNext(func () bool{return true},builderE)

	mergeBuilder.SetNext(func () bool{return true},builderE)
	error:=stateMachine.AddBuilder(builderA)
	if error!=nil{
		t.Fatalf("error: %s",error)
	}
	for i:=0;i<5;i++{
		e:=stateMachine.Clock()
		if e!=nil{
			if e.Error()!="invalid transition: cannot wait for a transition that doesn't merge back"{
				t.Fatalf("expected: invalid transition: cannot wait for a transition that doesn't merge back, got: %s",e)
			}
		}
	}
}
func StateMachineWithDualPathTakeEandB(res *string,takeB bool) StateMachine{
	stateMachine := CreateStateMachine()

	builderA:=CreateBuilderStateBase("a")
	builderA.SetActionDo(func() error {
		*res+="doA->"
		return nil
	})
	builderA.SetExitAction(func() error {
		*res+="exitA->"
		return nil
	})
	builderA.SetEntryAction(func() error {
		*res+="entryA->"
		return nil
	})
	
	builderB:=CreateBuilderStateBase("b")
	builderB.SetActionDo(func() error {
		*res+="doB->"
		return nil
	})
	builderB.SetExitAction(func() error {
		*res+="exitB->"
		return nil
	})
	builderB.SetEntryAction(func() error {
		*res+="entryB->"
		return nil
	})
	
	builderE:=CreateBuilderStateBase("e")
	
	builderE.SetEntryAction(func() error {
		*res+="entryE->"
		return nil
	})
	
	builderE.SetExitAction(func() error {
		*res+="exitE->"
		return nil
	})
	builderE.SetActionDo(func() error {
		*res+="doE->"
		return nil
	})
	builderA.AddBranch(func () bool{return takeB},builderB)
	builderA.AddBranch(func () bool{return !takeB},builderE)

	stateMachine.AddBuilder(builderA)
	return *stateMachine
}
func TestStateMachineWithDualPathTakeE(t *testing.T) {
	res:=""
	stateMachine:=StateMachineWithDualPathTakeEandB(&res,false)
	stateMachine.Clock()
	stateMachine.Clock()
	stateMachine.Clock()
	if res!="entryA->exitA->entryE->exitE->"{
		t.Fatalf("expected: entryA->exitA->entryE->exitE->, got: %s",res)
	}
}
func TestStateMachineWithDualPathTakeB(t *testing.T) {
	res:=""
	stateMachine:=StateMachineWithDualPathTakeEandB(&res,true)
	stateMachine.Clock()
	stateMachine.Clock()
	stateMachine.Clock()
	if res!="entryA->exitA->entryB->exitB->"{
		t.Fatalf("expected: entryA->doA->exitA->entryB->exitB->, got: %s",res)
	}
}
func TestCompositeState(t *testing.T) {
	res:=""
	stateA:=CreateBuilderStateBase("a")
	stateA.SetActionDo(func() error {
		res+="doA->"
		return nil
	})
	stateA.SetExitAction(func() error {
		res+="exitA->"
		return nil
	})
	stateA.SetEntryAction(func() error {
		res+="entryA->"
		return nil
	})
	stateB:=CreateBuilderStateBase("b")
	stateB.SetEntryAction(func() error {
		res+="entryB->"
		return nil
	})
	stateB.SetActionDo(func() error {
		res+="doB->"
		return nil
	})
	stateB.SetExitAction(func() error {
		res+="exitB->"
		return nil
	})
	stateC:=CreateBuilderStateBase("c")
	stateC.SetEntryAction(func() error {
		res+="entryC->"
		return nil
	})
	stateC.SetActionDo(func() error {
		res+="doC->"
		return nil
	})
	stateC.SetExitAction(func() error {
		res+="exitC->"
		return nil
	})
	stateD:=CreateBuilderStateBase("D")
	stateD.SetEntryAction(func() error {
		res+="entryD->"
		return nil
	})
	stateD.SetActionDo(func() error {
		res+="doD->"
		return nil
	})
	stateD.SetExitAction(func() error {
		res+="exitD->"
		return nil
	})

	composite:=CreateBuilderStateComposite("prova")
	composite.SetEntryAction(func() error {
		res+="entryComp->"
		return nil
	})
	composite.SetExitAction(func() error {
		res+="exitComp->"
		return nil
	})
	composite.SetActionDo(func() error {
		res+="doComp->"
		return nil
	})
	composite.AddState(stateB)
	composite.AddState(stateC)

	stateA.AddBranch(func () bool{return true},stateB)
	stateB.AddBranch(func () bool{return true},stateC)
	stateC.AddBranch(func () bool{return true},stateD)
	stateD.AddBranch(func () bool{return true},composite)

	stateMachine:=CreateStateMachine()
	stateMachine.AddBuilder(stateA)
	stateMachine.AddBuilderComposite(composite)
	stateMachine.Clock()
	stateMachine.Clock()
	stateMachine.Clock()
	stateMachine.Clock()
	stateMachine.Clock()
	stateMachine.Clock()

	if res!="entryA->exitA->entryComp->entryB->exitB->entryC->exitC->exitComp->entryD->exitD->entryComp->entryB->exitB->entryC->exitC->exitComp->entryD->"{
		t.Fatalf("entryA->exitA->entryComp->entryB->exitB->entryC->exitC->exitComp->entryD->exitD->entryComp->entryB->exitB->entryC->exitC->exitComp->entryD->, but got %s ",res)
	}

}


