package main

// callback for event on the carosello, the parameter is the index of the element
type CallBackCarosello func(int)
type CaroselloElement struct {
	wakeUpCallBack func(int)
	sleepCallBack  func(int)
	updateCallBack func(int)
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
	if len(e.elements) >= 0 {
		e.elements[0].updateCallBack(0)
		e.elements[0].wakeUpCallBack(0)
	}
	for i := 1; i < e.limitElements && i<len(e.elements); i++ {
		e.elements[i].updateCallBack(i)
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
	startRange:=e.index-e.index%e.limitElements
	if startRange<0{
		startRange=0
	}
	for i:=startRange;i<e.index+e.limitElements - e.index%e.limitElements;i++{
		if i>=len(e.elements){
			break
		}
		e.elements[i].sleepCallBack(i)
		if updateElement{
			e.elements[i].updateCallBack(i)
		}
	}
	e.elements[e.index].wakeUpCallBack(e.index)
}
