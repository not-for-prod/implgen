package test

import (
	_ "embed"
)

type Test struct {
}

func NewTest() *Test {
	return &Test{}
}
