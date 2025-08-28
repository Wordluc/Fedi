package main

import (
	"github.com/Wordluc/GTUI/Core"
	"github.com/Wordluc/GTUI/Core/Component"
)

type Callbacks[tDisplay Core.IContainer, tdata any] struct {
	updateDisplay   func(display tDisplay, data tdata)
	newDisplay      func(nDisplay int) tDisplay
	selectDisplay   func(nDisplay tDisplay)
	deselectDisplay func(nDisplay tDisplay)
}
type Carosello[tDisplay Core.IContainer, tdata any] struct {
	displayElements        []tDisplay
	currentDisplaySelected int
	currentDataSelected    int
	dataElements           []tdata
	callback               Callbacks[tDisplay, tdata]
	*Component.Container
}

func CreateCarosello[tDisplay Core.IContainer, tdata any](nDisplayElements int, callback Callbacks[tDisplay, tdata]) *Carosello[tDisplay, tdata] {
	container := Component.CreateContainer(0, 0)
	displays := make([]tDisplay, nDisplayElements)
	for i := range nDisplayElements {
		ele := callback.newDisplay(i)
		displays[i] = ele
		callback.deselectDisplay(ele)
		container.AddContainer(ele)
	}
	container.SetLayer(1)
	return &Carosello[tDisplay, tdata]{
		displayElements: displays,
		dataElements:    make([]tdata, 0),
		callback:        callback,
		Container:       container,
	}
}

func (c *Carosello[tDisplay, tdata]) AddData(data tdata) {
	c.dataElements = append(c.dataElements, data)
	if len(c.dataElements) > 0 {
		for i := range c.displayElements {
			c.callback.updateDisplay(c.displayElements[i], c.dataElements[i%len(c.dataElements)])
		}
	}
	if len(c.dataElements) == 0 {
		return
	}
	c.callback.selectDisplay(c.displayElements[c.currentDisplaySelected])
}

func (c *Carosello[tDisplay, tdata]) DeleteData(i int) {
	c.callback.deselectDisplay(c.displayElements[c.currentDisplaySelected])
	if i < 0 || i >= len(c.dataElements) {
		return
	}
	c.dataElements = append(c.dataElements[:i], c.dataElements[i+1:]...)
	c.currentDataSelected--
	c.currentDisplaySelected--
	if c.currentDisplaySelected < 0 {
		c.currentDisplaySelected = len(c.displayElements) - 1
	}
	if c.currentDataSelected < 0 {
		c.currentDataSelected = len(c.dataElements) - 1
	}
	c.updateDisplay()
	c.callback.selectDisplay(c.displayElements[c.currentDisplaySelected])

}
func (c *Carosello[tDisplay, tdata]) Next() {
	c.callback.deselectDisplay(c.displayElements[c.currentDisplaySelected])
	if len(c.dataElements) == 0 {
		return
	}
	c.currentDisplaySelected++
	c.currentDataSelected++
	if c.currentDisplaySelected >= len(c.displayElements) {
		c.currentDisplaySelected = 0
	}
	if c.currentDataSelected >= len(c.dataElements) {
		c.currentDataSelected = 0
	}
	c.callback.selectDisplay(c.displayElements[c.currentDisplaySelected])
	if c.currentDisplaySelected%len(c.displayElements) == 0 {
		c.updateDisplay()
	}
}
func (c *Carosello[tDisplay, tdata]) Pre() {
	c.callback.deselectDisplay(c.displayElements[c.currentDisplaySelected])
	if len(c.dataElements) == 0 {
		return
	}
	c.currentDisplaySelected--
	c.currentDataSelected--
	if c.currentDisplaySelected < 0 {
		c.currentDisplaySelected = len(c.displayElements) - 1
	}
	if c.currentDataSelected < 0 {
		c.currentDataSelected = len(c.dataElements) - 1
	}
	c.callback.selectDisplay(c.displayElements[c.currentDisplaySelected])
	if c.currentDisplaySelected == len(c.displayElements)-1 {
		c.updateDisplay()
	}
}

func (c *Carosello[tDisplay, tdata]) GetSelectedElement() (int, tdata) {
	if len(c.dataElements) == 0 {
		return 0, *new(tdata)
	}
	return c.currentDataSelected, c.dataElements[c.currentDataSelected]
}

func (c *Carosello[tDisplay, tdata]) updateDisplay() {
	for iDisplay := range c.displayElements {
		if len(c.dataElements) == 0 {
			c.callback.updateDisplay(c.displayElements[iDisplay], *new(tdata))
			continue
		}
		data := c.dataElements[(c.currentDataSelected+iDisplay)%len(c.dataElements)]
		c.callback.updateDisplay(c.displayElements[iDisplay], data)
	}
}
