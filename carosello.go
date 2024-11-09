package main

// callback for event on the carosello, the parameter is the index of the element
type CallBackCarosello func(int)
type CaroselloElement struct {
	index int
	wakeUpCallBack func()
	sleepCallBack  func()
	updateCallBack func()
}

type Carosello struct {
	index         int
	elements      []*CaroselloElement
	limitElements int
}

func CreateCarosello(x, y int, limit int) *Carosello {
	return &Carosello{
		index:         0,
		elements:      make([]*CaroselloElement, 0),
		limitElements: limit,
	}
}

func (e *Carosello) AddElement(element *CaroselloElement) {
	e.elements = append(e.elements, element)
	for i := 0; i < e.limitElements && i<len(e.elements); i++ {
		e.elements[i].updateCallBack()
	}
}
func (e *Carosello) NextOrPre(isPre bool) {
	pre_relativeIndex := e.index % e.limitElements
	if isPre {
		e.index--
	} else {
		e.index++
	}
	if e.index < 0 {
		e.index = len(e.elements) - 1
		e.updateElement(true)
		return
	}
	if e.index == len(e.elements) {
		e.index = 0
		e.updateElement(true)
		return
	}
	post_relativeIndex := e.index % e.limitElements
	if post_relativeIndex == 0 && pre_relativeIndex == e.limitElements-1 {
		e.updateElement(true)
	}
	if post_relativeIndex == e.limitElements-1 && pre_relativeIndex == 0 {
		e.updateElement(true)
	}
	e.updateElement(false)
}
func (e *Carosello) updateElement(updateElement bool) {
	if len(e.elements) == 0 {
		return
	}
	startRange:=e.index-e.index%e.limitElements
	if startRange<0{
		startRange=0
	}
	for i:=startRange;i<e.index+e.limitElements - e.index%e.limitElements;i++{
		if i>=len(e.elements){
			break
		}
		e.elements[i].sleepCallBack()
		if updateElement{
			e.elements[i].updateCallBack()
		}
	}
	e.elements[e.index].wakeUpCallBack()
}

func (e *Carosello) ForEachElements(action func (*CaroselloElement)) {
	for i := 0; i < len(e.elements); i++ {
		action(e.elements[i])
	}
}
