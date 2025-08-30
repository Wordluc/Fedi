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
	firstElementInPage     int
	nPages                 int
	dataElements           []tdata
	callback               Callbacks[tDisplay, tdata]
	*Component.Container
}

func CreateCarosello[tDisplay Core.IContainer, tdata any](nDisplayElements int, callback Callbacks[tDisplay, tdata]) *Carosello[tDisplay, tdata] {
	container := Component.CreateContainer()
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

func (c *Carosello[tDisplay, tdata]) AddDataAll(data ...tdata) {
	for i := range data {
		c.AddData(data[i])
	}
}

func (c *Carosello[tDisplay, tdata]) AddData(data tdata) {
	c.dataElements = append(c.dataElements, data)
	if len(c.dataElements) > 0 {
		for i := range c.displayElements {
			c.callback.updateDisplay(c.displayElements[i], c.dataElements[i%len(c.dataElements)])
		}
	}
	c.nPages = (len(c.dataElements) / len(c.displayElements))
	c.callback.selectDisplay(c.displayElements[c.currentDisplaySelected])
}

func (c *Carosello[tDisplay, tdata]) Refresh(data ...tdata) {
	if len(c.dataElements) == len(data) {
		c.dataElements = nil
		c.dataElements = append(c.dataElements, data...)
		c.updateDisplays()
	}
}
func (c *Carosello[tDisplay, tdata]) Reset() {
	c.callback.deselectDisplay(c.displayElements[c.currentDisplaySelected])
	c.nPages = 0
	c.currentDisplaySelected = 0
	c.firstElementInPage = 0
	c.dataElements = nil
	c.updateDisplays()
}

func (c *Carosello[tDisplay, tdata]) DeleteData(i int) {
	c.callback.deselectDisplay(c.displayElements[c.currentDisplaySelected])
	if i < 0 || i >= len(c.dataElements) {
		return
	}
	c.dataElements = append(c.dataElements[:i], c.dataElements[i+1:]...)
	c.currentDisplaySelected--
	if c.currentDisplaySelected < 0 {
		c.currentDisplaySelected = 0
		c.firstElementInPage = (c.firstElementInPage - 1)
	}
	c.updateDisplays()
	if len(c.dataElements) != 0 {
		c.callback.selectDisplay(c.displayElements[c.currentDisplaySelected])
	}
}

func (c *Carosello[tDisplay, tdata]) Next() {
	c.callback.deselectDisplay(c.displayElements[c.currentDisplaySelected])
	if len(c.dataElements) == 0 {
		return
	}
	c.currentDisplaySelected++
	if c.currentDisplaySelected >= len(c.displayElements) {
		c.currentDisplaySelected = len(c.displayElements) - 1
		c.firstElementInPage = (c.firstElementInPage + 1)
	}

	c.callback.selectDisplay(c.displayElements[c.currentDisplaySelected])
	if c.currentDisplaySelected == len(c.displayElements)-1 {
		c.updateDisplays()
	}
}
func (c *Carosello[tDisplay, tdata]) Pre() {
	c.callback.deselectDisplay(c.displayElements[c.currentDisplaySelected])
	if len(c.dataElements) == 0 {
		return
	}
	c.currentDisplaySelected--
	if c.currentDisplaySelected < 0 {
		c.currentDisplaySelected = 0
		c.firstElementInPage = (c.firstElementInPage - 1)
		if c.firstElementInPage < 0 {
			c.firstElementInPage = len(c.dataElements) - 1
		}
	}
	c.callback.selectDisplay(c.displayElements[c.currentDisplaySelected])
	if c.currentDisplaySelected == 0 {
		c.updateDisplays()
	}
}

func (c *Carosello[tDisplay, tdata]) GetSelectedElement() (int, tdata) {
	if len(c.dataElements) == 0 {
		return 0, *new(tdata)
	}
	i := (c.firstElementInPage + c.currentDisplaySelected) % len(c.dataElements)
	data := c.dataElements[i]
	return i, data
}

func (c *Carosello[tDisplay, tdata]) GetElements() []tdata {
	return c.dataElements
}

func (c *Carosello[tDisplay, tdata]) updateDisplays() {
	for iDisplay := range c.displayElements {
		if len(c.dataElements) == 0 {
			c.callback.updateDisplay(c.displayElements[iDisplay], *new(tdata))
			continue
		}
		data := c.dataElements[(c.firstElementInPage+iDisplay)%len(c.dataElements)]
		c.callback.updateDisplay(c.displayElements[iDisplay], data)
	}
}
