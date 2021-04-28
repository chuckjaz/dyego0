package sources

import (
	"dyego0/ast"
	"dyego0/tokens"
)

// Source is the result of parsering the text from a SourceReader
type Source interface {
	// Text returns the text of the source from start to end
	Text(start, end int) string

	// File is the tokens File for the source
	File() tokens.File

	// Ast is the AST of the source
	Ast() ast.Element
}

// SourceSet is a set of related source files
type SourceSet interface {
	// FileSet the tokens FileSet to map locations back to files
	FileSet() tokens.FileSet

	// Source returns a Source from the source set by filename
	Source(filename string) Source

	// SourceOf returns a Source from teh source set by tokens File
	SourceOf(file tokens.File) Source
}
