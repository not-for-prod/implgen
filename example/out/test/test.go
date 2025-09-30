package test

import (
	_ "embed"
)

type Test struct {
}

func NewImplementation() *Test {
	return &Test{}
}
