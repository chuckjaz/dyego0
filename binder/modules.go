package binder

import (
	"io"
)

// ModuleSourceScope is a scope for finding module sources
type ModuleSourceScope interface {
	// FindScope finds a subscope of a module scope
	FindScope(name string) (ModuleSourceScope, error)

	// Find a source module in a scope
	Find(name string) (ModuleSource, error)
}

// ModuleSource is a source file for a module
type ModuleSource interface {
	// Name is the name the module that appears in the scope
	Name() string

	// FileName is the file name use in  error messages
	FileName() string

	// NewReader create a reader for the source
	NewReader() (io.Reader, error)
}

// NewModuleSource creates a new ModuleSource
func NewModuleSource(name string, fileName string, readFactory func() (io.Reader, error)) ModuleSource {
	return &moduleSource{name: name, fileName: fileName, readFactory: readFactory}
}

type moduleSource struct {
	name        string
	fileName    string
	readFactory func() (io.Reader, error)
}

func (m *moduleSource) Name() string {
	return m.name
}

func (m *moduleSource) FileName() string {
	return m.fileName
}

func (m *moduleSource) NewReader() (io.Reader, error) {
	return m.readFactory()
}
