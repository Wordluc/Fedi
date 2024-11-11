package main

//import "github.com/Wordluc/GTUI/Core/Utils"

// callback for event on the carosello, the parameter is the index of the element
type CallBackCarosello func(int)
type CaroselloElement struct {
	index          int
	wakeUpCallBack func(int)
	sleepCallBack  func(int)
	updateCallBack func(int)
}

type Carosello struct {
	index             int
	startRangeElement int
	selectedElement   int
	elements          []*CaroselloElement
	limitElements     int
}

func CreateCarosello(x, y int, limit int) *Carosello {
	return &Carosello{
		index:             0,
		startRangeElement: 0,
		selectedElement:   0,
		elements:          make([]*CaroselloElement, 0),
		limitElements:     limit,
	}
}

func (e *Carosello) AddElement(element *CaroselloElement) {
	e.elements = append(e.elements, element)
	i := len(e.elements) - 1
	if i < e.limitElements {
		e.elements[i].updateCallBack(i % e.limitElements)
	}
}
func (e *Carosello) NextOrPre(isPre bool) {
	pre_relativeIndex := e.selectedElement
	if isPre {
		e.index--
		e.selectedElement--
	} else {
		e.index++
		e.selectedElement++
	}

	if e.selectedElement == -1 {
		e.selectedElement = e.limitElements - 1
	}
	if e.selectedElement == e.limitElements {
		e.selectedElement = 0
	}
	isGoingDown := e.selectedElement == 0 && pre_relativeIndex == e.limitElements-1
	isGoingUp := e.selectedElement == e.limitElements-1 && pre_relativeIndex == 0
	if isGoingDown { //verso basso
		e.startRangeElement = e.startRangeElement + e.limitElements
		if e.startRangeElement >= len(e.elements) {
			e.startRangeElement = e.startRangeElement - len(e.elements)
		}
		e.updateElement(true)
		return
	}
	if isGoingUp { //verso il l'alto
		e.startRangeElement = e.startRangeElement - e.limitElements
		if e.startRangeElement < 0 {
			e.startRangeElement = len(e.elements) + e.startRangeElement
		}
		e.updateElement(true)
		return
	}
	e.updateElement(false)
}
func (e *Carosello) updateElement(updateElement bool) {
	iblock := 0
	for i := e.startRangeElement; i < e.startRangeElement+e.limitElements; i++ {
		e.elements[i%len(e.elements)].sleepCallBack(iblock)
		if updateElement {
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
	e.elements[e.index%len(e.elements)].wakeUpCallBack(e.selectedElement)

}

func (e *Carosello) ForEachElements(action func(*CaroselloElement, int)) {
	for i := 0; i < len(e.elements); i++ {
		action(e.elements[i], i%e.limitElements)
	}
}
