package sources

import (
	"fmt"
	"io"
	"path"
)

// NewFilesSourceReaderScope creates a file path file module source scope
func NewSourceFileReaderScope(files []string, reader func(filename string) (io.Reader, error)) (SourceReaderScope, error) {
	result := newFilesScope()
	for _, file := range files {
		names := splitName(noExt(file))
		if len(names) > 1 {
			nameIndex := len(names) - 1
			err := result.populate(file, reader, names[0:nameIndex], names[nameIndex])
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Invalid file name '%s'", file)
		}
	}
	return result, nil
}

type fileReader struct {
	name     string
	filename string
	reader   func(filename string) (io.Reader, error)
}

func (f *fileReader) Name() string {
	return f.name
}

func (f *fileReader) FileName() string {
	return f.filename
}

func (f *fileReader) NewReader() (io.Reader, error) {
	file, err := f.reader(f.filename)
	return file, err
}

type filesScope struct {
	files  map[string]*fileReader
	scopes map[string]*filesScope
}

func (s *filesScope) Find(name string) (SourceReader, error) {
	file, ok := s.files[name]
	if !ok {
		return nil, fmt.Errorf("File '%s' not found", name)
	}
	return file, nil
}

func (s *filesScope) FindScope(name string) (SourceReaderScope, error) {
	scope, ok := s.scopes[name]
	if !ok {
		return nil, fmt.Errorf("Directory '%s' not found", name)
	}
	return scope, nil
}

func newFilesScope() *filesScope {
	return &filesScope{files: make(map[string]*fileReader), scopes: make(map[string]*filesScope)}
}

func (s *filesScope) populate(filename string, reader func(filename string) (io.Reader, error), dir []string, name string) error {
	if len(dir) == 0 {
		_, ok := s.files[name]
		if ok {
			return fmt.Errorf("Duplicate file '%s'", filename)
		}
		s.files[name] = &fileReader{name: name, filename: filename, reader: reader}
		return nil
	}
	scopeName := dir[0]
	scope, ok := s.scopes[scopeName]
	if !ok {
		scope = newFilesScope()
		s.scopes[scopeName] = scope
	}
	return scope.populate(filename, reader, dir[1:], name)
}

func splitName(filename string) []string {
	var result []string

	current := filename
	for current != "/" && current != "." && current != "" {
		base, name := path.Split(current)
		var empty []string
		result = append(append(empty, name), result...)
		if len(base) > 1 {
			current = base[0 : len(base)-1]
		} else {
			current = base
		}
	}
	return result
}

func noExt(filename string) string {
	ext := path.Ext(filename)
	return filename[0 : len(filename)-len(ext)]
}
