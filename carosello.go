package main

import "fmt"

type CallBackCarosello func(int)
type CaroselloElement struct {
	//Is called when the element is woken up, passing the index of the block as an argument
	wakeUpCallBack func(int)
	//Is called when the element is put to sleep, passing the index of the block as an argument
	sleepCallBack func(int)
	//Is called when the element needs to be updated, passing the index of the block as an argument
	updateCallBack func(int)
}

type Carosello struct {
	index             int
	startRangeElement int
	selectedBlock     int
	elements          []*CaroselloElement
	limitBlocks       int
}

func CreateCarosello(x, y int, limitBlock int) *Carosello {
	return &Carosello{
		index:             0,
		startRangeElement: 0,
		elements:          make([]*CaroselloElement, 0),
		selectedBlock:     0,
		limitBlocks:       limitBlock,
	}
}

func (e *Carosello) AddElement(element *CaroselloElement) {
	e.elements = append(e.elements, element)
	i := len(e.elements) - 1
	if i < e.limitBlocks {
		e.elements[i].updateCallBack(i % e.limitBlocks)
	}
}
func (e *Carosello) NextOrPre(isPre bool) {
	if len(e.elements)==0{
		return
	}
	pre_selectedBlock := e.selectedBlock
	if isPre {
		e.index--
		e.selectedBlock--
	} else {
		e.index++
		e.selectedBlock++
	}

	if e.selectedBlock == -1 {
		e.selectedBlock = e.limitBlocks - 1
	}
	if e.selectedBlock == e.limitBlocks {
		e.selectedBlock = 0
	}
	isGoingDown := e.selectedBlock == 0 && pre_selectedBlock == e.limitBlocks-1
	isGoingUp := e.selectedBlock == e.limitBlocks-1 && pre_selectedBlock == 0
	if isGoingDown {
		e.startRangeElement = e.startRangeElement + e.limitBlocks
		if e.startRangeElement >= len(e.elements) {
			e.startRangeElement = e.startRangeElement - len(e.elements)
		}
		e.UpdateElementState(true,true)
		return
	}
	if isGoingUp {
		e.startRangeElement = e.startRangeElement - e.limitBlocks
		for e.startRangeElement < 0 {
			e.startRangeElement = len(e.elements) + e.startRangeElement
		}
		e.UpdateElementState(true,true)
		return
	}
	e.UpdateElementState(false,true)
}
func (e *Carosello) UpdateElementState(refreshContentElement bool, setWakeup bool) {
	iblock := 0
	if len(e.elements)==0{
		return
	}
	for i := e.startRangeElement; i < e.startRangeElement+e.limitBlocks; i++ {
		e.elements[i%len(e.elements)].sleepCallBack(iblock)
		if refreshContentElement {
			e.elements[i%len(e.elements)].updateCallBack(iblock)
		}
		iblock++
	}
	if e.index == len(e.elements) {
		e.index = 0
	}
	if e.index == -1 {
		e.index = len(e.elements) - 1
	}
	if setWakeup {
		e.elements[e.index%len(e.elements)].wakeUpCallBack(e.selectedBlock)
	}
}
func (e *Carosello) SetIndex(iNeeded int) error {
	i := e.index
	startI := e.index
	for {
		if i == iNeeded {
			return nil
		}
		e.NextOrPre(i < iNeeded)
		i = e.index
		if startI == i {
			return fmt.Errorf("index not found")
		}
	}
}
func (e *Carosello) SleepAll() {
	iblock := 0
	if len(e.elements)==0{
		return
	}
	for i := e.startRangeElement; i < e.startRangeElement+e.limitBlocks; i++ {
		e.elements[i%len(e.elements)].sleepCallBack(iblock)
		iblock++
	}
}
func (e *Carosello) ForEachElements(action func(*CaroselloElement, int)) {
	if len(e.elements)==0{
		return
	}
	for i := 0; i < len(e.elements); i++ {
		action(e.elements[i], i%e.limitBlocks)
	}
}
func (e *Carosello) GetIntex() int {
	return e.index
}
func (e *Carosello) GetSelected() int {
	return e.selectedBlock
}
func (e *Carosello) GetElementsNumber() (i int) {
	return len(e.elements)
}
