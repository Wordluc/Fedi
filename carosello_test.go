package main

import (
	"fmt"
	"strings"
	"testing"
)

func createCaroselloElementHelperMethods(iEl int, act *string,showOnlyUpdate bool) *CaroselloElement {
	return &CaroselloElement{
		index: iEl,
		wakeUpCallBack: func(i int) {
			if showOnlyUpdate {
				return
			}
			(*act) = (*act) + "posBlock:"+ fmt.Sprint(i) + ",wakeUpCallBack,pos"+fmt.Sprint(iEl)+"\n"
		},
		sleepCallBack: func(i int) {
			if showOnlyUpdate {
				return
			}
			(*act) = (*act) + "posBlock:"+ fmt.Sprint(i) + ",sleepCallCallBack,pos"+fmt.Sprint(iEl)+"\n"
		},
		updateCallBack: func(i int) {
			(*act) = (*act) + "*posBlock:" + fmt.Sprint(i) + ",updateCallCallBack,pos" + fmt.Sprint(iEl) + "\n"
		},
	}
}
func checkresult(act string, exp string) bool {
	exp=strings.Replace(exp, "\n", "",-1 )
	exp=strings.Replace(exp, "	", "",-1 )
	act=strings.Replace(act, "\n", "",-1 )
	act=strings.Replace(act, "	", "",-1 )
	return act == exp
}
func TestCaroselloChangePageNextFromDown(t *testing.T) {
	carosello := CreateCarosello(0, 0, 3)
	var result *string = new(string)
	for i := 0; i < 5; i++ {
		carosello.AddElement(createCaroselloElementHelperMethods(i, result,true))
	}
	carosello.updateElement(false)
	(*result) += "-----------------\n"
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	exp := `
		*posBlock:0,updateCallCallBack,pos0
		*posBlock:1,updateCallCallBack,pos1
		*posBlock:2,updateCallCallBack,pos2
		-----------------
		*posBlock:0,updateCallCallBack,pos3
		*posBlock:1,updateCallCallBack,pos4
		*posBlock:2,updateCallCallBack,pos0`
	if checkresult(*result, exp) == false {
		t.Errorf("expected \n%s\n got\n%s\n", exp, *result)
	}
}
func TestCaroselloChange2PagesNextFromDown(t *testing.T) {
	carosello := CreateCarosello(0, 0, 3)
	var result *string = new(string)
	for i := 0; i < 5; i++ {
		carosello.AddElement(createCaroselloElementHelperMethods(i, result,true))
	}
	carosello.updateElement(false)
	(*result) += "-----------------\n"
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	(*result) += "-----------------\n"
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	exp := `
*posBlock:0,updateCallCallBack,pos0
*posBlock:1,updateCallCallBack,pos1
*posBlock:2,updateCallCallBack,pos2
-----------------
*posBlock:0,updateCallCallBack,pos3
*posBlock:1,updateCallCallBack,pos4
*posBlock:2,updateCallCallBack,pos0
-----------------
*posBlock:0,updateCallCallBack,pos1
*posBlock:1,updateCallCallBack,pos2
*posBlock:2,updateCallCallBack,pos3
	`
	if checkresult(*result, exp) == false {
		t.Errorf("expected \n%s\n got\n%s\n", exp, *result)
	}
}
func TestCaroselloChange3PagesNextFromDown(t *testing.T) {
	carosello := CreateCarosello(0, 0, 3)
	var result *string = new(string)
	for i := 0; i < 5; i++ {
		carosello.AddElement(createCaroselloElementHelperMethods(i, result,true))
	}
	carosello.updateElement(false)
	(*result) += "-----------------\n"
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	(*result) += "-----------------\n"
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	(*result) += "-----------------\n"
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	carosello.NextOrPre(false)
	exp := `
*posBlock:0,updateCallCallBack,pos0
*posBlock:1,updateCallCallBack,pos1
*posBlock:2,updateCallCallBack,pos2
-----------------
*posBlock:0,updateCallCallBack,pos3
*posBlock:1,updateCallCallBack,pos4
*posBlock:2,updateCallCallBack,pos0
-----------------
*posBlock:0,updateCallCallBack,pos1
*posBlock:1,updateCallCallBack,pos2
*posBlock:2,updateCallCallBack,pos3
-----------------
*posBlock:0,updateCallCallBack,pos4
*posBlock:1,updateCallCallBack,pos0
*posBlock:2,updateCallCallBack,pos1
	`
	if checkresult(*result, exp) == false {
		t.Errorf("expected \n%s\n got\n%s\n", exp, *result)
	}
}
func TestCaroselloChangePagesNextFromUp(t *testing.T) {
	carosello := CreateCarosello(0, 0, 3)
	var result *string = new(string)
	for i := 0; i < 5; i++ {
		carosello.AddElement(createCaroselloElementHelperMethods(i, result,true))
	}
	(*result) += "-----------------\n"
	carosello.NextOrPre(true)
	carosello.NextOrPre(true)
	carosello.NextOrPre(true)
	(*result) += "-----------------\n"
	carosello.NextOrPre(true)
	carosello.NextOrPre(true)
	carosello.NextOrPre(true)
	(*result) += "-----------------\n"
	carosello.NextOrPre(true)
	carosello.NextOrPre(true)
	carosello.NextOrPre(true)
	exp := `
*posBlock:0,updateCallCallBack,pos0
*posBlock:1,updateCallCallBack,pos1
*posBlock:2,updateCallCallBack,pos2
-----------------
*posBlock:0,updateCallCallBack,pos2
*posBlock:1,updateCallCallBack,pos3
*posBlock:2,updateCallCallBack,pos4
-----------------
*posBlock:0,updateCallCallBack,pos4
*posBlock:1,updateCallCallBack,pos0
*posBlock:2,updateCallCallBack,pos1
-----------------
*posBlock:0,updateCallCallBack,pos1
*posBlock:1,updateCallCallBack,pos2
*posBlock:2,updateCallCallBack,pos3

	`
	if checkresult(*result, exp) == false {
		t.Errorf("expected \n%s\n got\n%s\n", exp, *result)
	}
}
