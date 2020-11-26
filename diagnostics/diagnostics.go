package diagnostics

import (
	"fmt"
	"strings"

	"dyego0/errors"
	"dyego0/tokens"
)

// SourceProvider is an interface that allows Format to find the content of a source file
type SourceProvider interface {
	// Source returns the Source for the given filename
	Source(filename string) Source
}

// Source is an interface that allows accessing the test of a source file
type Source interface {
	// Text returns the text of the source from start to end
	Text(start, end int) string
}

// Format formats error messages into a string. If a sourceProvider is provided then souce line containing
// error is added as well.
func Format(errors []errors.Error, fileSet tokens.FileSet, sourceProvider SourceProvider) string {
	var result string
	if errors != nil {
		for _, err := range errors {
			start := err.Start()
			position := fileSet.Position(start)
			result += fmt.Sprintf("%s: %s\n", position, err.Error())
			if sourceProvider != nil {
				file := fileSet.File(start)
				if file != nil {
					source := sourceProvider.Source(position.FileName())
					if source != nil {
						lineStart := file.LineStart(position.Line())
						lineEnd := file.LineStart(position.Line() + 1)
						text := source.Text(lineStart, lineEnd)
						if len(text) > 0 && text[len(text)-1] == '\n' {
							text = text[0 : len(text)-1]
						}
						end := err.End()
						endPosition := file.Position(end)
						startColumn := position.Column()
						endColumn := endPosition.Column()
						if position.Line() != endPosition.Line() {
							endColumn = len(text) + 1
						}
						result += text + "\n"
						if startColumn > 0 {
							for _, ch := range text[0 : startColumn-1] {
								if ch == '\t' {
									result += "\t"
								} else {
									result += " "
								}
							}
						}
						result += strings.Repeat("^", endColumn-startColumn) + "\n"
					}
				}
			}
		}
	}
	return result
}
