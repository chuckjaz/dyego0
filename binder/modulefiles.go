package binder

import (
	"fmt"
	"io"
	"path"
)

// NewFilesModuleSourceScope creates a file path file module source scope
func NewFilesModuleSourceScope(files []string, reader func(fileName string) (io.Reader, error)) (ModuleSourceScope, error) {
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

type fileModule struct {
	name     string
	fileName string
	reader   func(fileName string) (io.Reader, error)
}

func (f *fileModule) Name() string {
	return f.name
}

func (f *fileModule) FileName() string {
	return f.fileName
}

func (f *fileModule) NewReader() (io.Reader, error) {
	file, err := f.reader(f.fileName)
	return file, err
}

type filesScope struct {
	files  map[string]*fileModule
	scopes map[string]*filesScope
}

func (s *filesScope) Find(name string) (ModuleSource, error) {
	file, ok := s.files[name]
	if !ok {
		return nil, fmt.Errorf("File '%s' not found", name)
	}
	return file, nil
}

func (s *filesScope) FindScope(name string) (ModuleSourceScope, error) {
	scope, ok := s.scopes[name]
	if !ok {
		return nil, fmt.Errorf("Directory '%s' not found", name)
	}
	return scope, nil
}

func newFilesScope() *filesScope {
	return &filesScope{files: make(map[string]*fileModule), scopes: make(map[string]*filesScope)}
}

func (s *filesScope) populate(fileName string, reader func(fileName string) (io.Reader, error), dir []string, name string) error {
	if len(dir) == 0 {
		_, ok := s.files[name]
		if ok {
			return fmt.Errorf("Duplicate file '%s'", fileName)
		}
		s.files[name] = &fileModule{name: name, fileName: fileName, reader: reader}
		return nil
	}
    scopeName := dir[0]
    scope, ok := s.scopes[scopeName]
    if !ok {
        scope = newFilesScope()
        s.scopes[scopeName] = scope
    }
     return scope.populate(fileName, reader, dir[1:], name)
}

func splitName(fileName string) []string {
	var result []string

	current := fileName
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

func noExt(fileName string) string {
	ext := path.Ext(fileName)
	return fileName[0 : len(fileName)-len(ext)]
}
