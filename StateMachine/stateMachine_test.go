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
	builderEnd:=CreateBuilderStateEnd("end")
	builderEnd.SetActionDo(func() error {
		res+="doEnd->"
		return nil
	})
	builderEnd.SetExitAction(func() error {
		res+="exitEnd->"
		return nil
	})
	builderEnd.SetEntryAction(func() error {
		res+="entryEnd->"
		return nil
	})

	builderA.SetNext(func () bool{return true},forkBuilder)
	forkBuilder.AddTo(func () bool{return true},builderC)
	forkBuilder.AddTo(func () bool{return true},builderB)

	builderC.SetNext(func () bool{return true},mergeBuilder)
	builderB.SetNext(func () bool{return true},mergeBuilder)

	mergeBuilder.AddToWait(builderC)
	mergeBuilder.AddToWait(builderB)

	mergeBuilder.SetNext(func () bool{return true},builderE)

	mergeBuilder.SetNext(func () bool{return true},builderE)
	builderE.SetNext(func () bool{return true},builderEnd)
	error:=stateMachine.AddBuilder(builderA)
	if error!=nil{
		t.Fatalf("error: %s",error)
	}
	for i:=0;i<6;i++{
		stateMachine.Clock()
	}
	exp:="entryA->doA->exitA->entryFork->exitFork->entryC->exitFork->entryB->doC->exitC->doB->exitB->entryE->entryE->exitE->entryEnd->exitE->entryEnd->doEnd->exitEnd->"
	if res!=exp{
		t.Fatalf("got: %s, expected: %s",res,exp)
	}
}
func TestStateMachineForkWithoutWait(t *testing.T) {
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
	builderEnd:=CreateBuilderStateEnd("end")
	builderEnd.SetActionDo(func() error {
		res+="doEnd->"
		return nil
	})
	builderEnd.SetExitAction(func() error {
		res+="exitEnd->"
		return nil
	})
	builderEnd.SetEntryAction(func() error {
		res+="entryEnd->"
		return nil
	})

	builderA.SetNext(func () bool{return true},forkBuilder)
	forkBuilder.AddTo(func () bool{return true},builderC)
	forkBuilder.AddTo(func () bool{return true},builderB)

	builderC.SetNext(func () bool{return true},mergeBuilder)
	builderB.SetNext(func () bool{return true},mergeBuilder)
//	mergeBuilder.AddToWait(builderC)
//	mergeBuilder.AddToWait(builderB)

	mergeBuilder.SetNext(func () bool{return true},builderE)

	mergeBuilder.SetNext(func () bool{return true},builderE)
	builderE.SetNext(func () bool{return true},builderEnd)
	error:=stateMachine.AddBuilder(builderA)
	if error!=nil{
		t.Fatalf("error: %s",error)
	}
	for i:=0;i<5;i++{
		e:=stateMachine.Clock()
		if e!=nil{
			t.Fatalf("error: %s",e)
		}
	}
	exp:="entryA->doA->exitA->entryFork->exitFork->entryC->exitFork->entryB->doC->exitC->entryE->doB->exitB->entryE->exitE->entryEnd->doEnd->exitEnd->"
	if res!=exp{
		t.Fatalf("got: %s, expected: |%s|",res,exp)
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
	builderEnd:=CreateBuilderStateEnd("end")
	builderEnd.SetActionDo(func() error {
		res+="doEnd->"
		return nil
	})
	builderEnd.SetExitAction(func() error {
		res+="exitEnd->"
		return nil
	})
	builderEnd.SetEntryAction(func() error {
		res+="entryEnd->"
		return nil
	})

	builderA.SetNext(func () bool{return true},forkBuilder)
	forkBuilder.AddTo(func () bool{return true},builderC)
	forkBuilder.AddTo(func () bool{return true},builderB)

	builderC.SetNext(func () bool{return true},mergeBuilder)
	builderB.SetNext(func () bool{return true},mergeBuilder)
	mergeBuilder.AddToWait(builderA)
//	mergeBuilder.AddToWait(builderB)

	mergeBuilder.SetNext(func () bool{return true},builderE)

	mergeBuilder.SetNext(func () bool{return true},builderE)
	builderE.SetNext(func () bool{return true},builderEnd)
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
