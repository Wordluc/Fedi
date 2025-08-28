package main

import (
	"os"
	"strings"
)

type Repositoty[t any] struct {
	fileName string
	parser   func(in []string) t
}

func NewRepositoty[t any](fileName string, parser func([]string) t) *Repositoty[t] {
	return &Repositoty[t]{
		fileName: fileName,
		parser:   parser,
	}
}

func (r *Repositoty[t]) Get() (res []t, err error) {
	file, e := os.ReadFile("prova.csv")
	if e != nil {
		return nil, e
	}
	lines := strings.Split(string(file), "\n")
	for i := range lines {
		if lines[i] == "" {
			continue
		}
		res = append(res, r.parser(strings.Split(lines[i], ";")))
	}
	return res, nil
}

func (r *Repositoty[t]) Add(input []string) error {
	inputCsv := strings.Join(input, ";")
	f, err := os.OpenFile(r.fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString("\n" + inputCsv)
	if err != nil {
		return err
	}
	return nil
}
