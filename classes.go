package ui

import "strings"

type Class struct {
	Class string
	Test  bool
}

type Classes []Class

func (classes Classes) String() string {
	var include []string
	for _, c := range classes {
		if c.Test {
			include = append(include, c.Class)
		}
	}

	return strings.Join(include, " ")
}
