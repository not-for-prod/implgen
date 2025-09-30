package model

// Package represents src interface with all package-specific attributes
type Package struct {
	// filepath.Base for go mod module name
	Name       string
	Interfaces []Interface
	Imports    []Import
}

type Import struct {
	Alias string
	Path  string
}

type Interface struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name     string
	In, Out  []Parameter
	Variadic *Parameter // optional
}

type Parameter struct {
	Name string
	Type string
}

type File struct {
	Path string
	Data []byte
}
