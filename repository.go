package main

import (
	"bufio"
	"os"
	"strings"
)

type Repositoty[t any] struct {
	fileName   string
	unmarshall func(in []string) t
	marshall   func(t) []string
	equal      func(t, t) bool
}

type TODO struct {
	Row    int
	Id     string
	Title  string
	Text   string
	Date   string
	Status MarkType
}

func NewRepositoty[t any](fileName string, unmarshall func([]string) t, marshall func(t) []string, equal func(t, t) bool) *Repositoty[t] {
	return &Repositoty[t]{
		fileName:   fileName,
		unmarshall: unmarshall,
		marshall:   marshall,
		equal:      equal,
	}
}

func (r *Repositoty[t]) Get() (res []t, err error) {
	f, err := os.OpenFile(r.fileName, os.O_RDONLY|os.O_CREATE, 00777)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		res = append(res, r.unmarshall(strings.Split(scanner.Text(), ";")))
	}
	return res, nil
}

func (r *Repositoty[t]) Add(input t) error {
	inputCsv := strings.Join(r.marshall(input), ";")
	f, err := os.OpenFile(r.fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + inputCsv)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repositoty[t]) Remove(ele t) error {
	content, e := r.Get()
	if e != nil {
		return e
	}
	os.Remove(r.fileName)
	for i := range content {
		if r.equal(content[i], ele) {
			continue
		}
		if err := r.Add(content[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repositoty[t]) Set(ele t) error {
	content, e := r.Get()
	if e != nil {
		return e
	}
	os.Remove(r.fileName)
	for i := range content {
		currentContent := content[i]
		if r.equal(content[i], ele) {
			currentContent = ele
		}
		if err := r.Add(currentContent); err != nil {
			return err
		}
	}
	return nil
}
