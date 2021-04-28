package sources

import (
	"io"
)

// SourceReaderScope is a scope for finding source readers
type SourceReaderScope interface {
	// FindScope finds a subscope of a module scope
	FindScope(name string) (SourceReaderScope, error)

	// Find a source module in a scope
	Find(name string) (SourceReader, error)
}

// SourceReader is a source reader
type SourceReader interface {
	// Name is the name the module that appears in the scope
	Name() string

	// FileName is the file name use in  error messages
	FileName() string

	// NewReader create a reader for the source
	NewReader() (io.Reader, error)
}

// NewModuleSource creates a new ModuleSource
func NewSourceReader(name string, fileName string, readFactory func() (io.Reader, error)) SourceReader {
	return &sourceReader{name: name, fileName: fileName, readFactory: readFactory}
}

type sourceReader struct {
	name        string
	fileName    string
	readFactory func() (io.Reader, error)
}

func (m *sourceReader) Name() string {
	return m.name
}

func (m *sourceReader) FileName() string {
	return m.fileName
}

func (m *sourceReader) NewReader() (io.Reader, error) {
	return m.readFactory()
}
