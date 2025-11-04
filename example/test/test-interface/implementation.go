package test_interface

import (
	_ "embed"
)

type Implementation struct {
}

func NewImplementation() *Implementation {
	return &Implementation{}
}
