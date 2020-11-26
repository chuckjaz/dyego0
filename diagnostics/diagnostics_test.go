package diagnostics

import (
	"dyego0/errors"
	"dyego0/location"
	"dyego0/tokens"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("diagnostics", func() {
	It("produces and empty string for an empty array of errors", func() {
		var errs []errors.Error
		msg := Format(errs, nil, nil)
		Expect(msg).To(Equal(""))
	})
	It("can format one error", func() {
		text := `Line one
			Line two
			Report location
			Line Four
		`
		fs, p, errs := buildErrors(text, "location", "This is the location")
		msg := Format(errs, fs, p)
		Expect(msg).To(Equal("file:3:11: This is the location\n\t\t\tReport location\n\t\t\t       ^^^^^^^^\n"))
	})
	It("can format multiple errors", func() {
		text := `Line one
			Line two
			Report location
			Line Four
			Report another location
			Line Six
			Line Seven
			Report a third location
			Line Nine
		`
		fs, p, errs := buildErrors(text, "location", "This is the location")
		msg := Format(errs, fs, p)
		Expect(strings.Contains(msg, "3:11")).To(BeTrue())
		Expect(strings.Contains(msg, "5:19")).To(BeTrue())
		Expect(strings.Contains(msg, "8:19")).To(BeTrue())
	})
	It("can format an error on the last line", func() {
		text := `Line one
			Line Nine
		report location`
		fs, p, errs := buildErrors(text, "location", "This is the location")
		msg := Format(errs, fs, p)
		Expect(msg).To(Equal("file:3:10: This is the location\n\t\treport location\n\t\t       ^^^^^^^^\n"))
	})
})

type sourceProvider struct {
	text string
}

func (p *sourceProvider) Source(name string) Source {
	if name == "file" {
		return p
	}
	return nil
}

func (p *sourceProvider) Text(start, end int) string {
	return p.text[start:end]
}

func providerFor(text string) SourceProvider {
	return &sourceProvider{text: text}
}

func fileFor(fs tokens.FileSet, text string) tokens.File {
	fb := fs.BuildFile("file", len(text))
	fb.AddLine(0)
	for offset, c := range text {
		if c == '\n' {
			fb.AddLine(offset + 1)
		}
	}
	return fb.Build()
}

func buildErrors(text, pattern, message string) (tokens.FileSet, SourceProvider, []errors.Error) {
	p := providerFor(text)
	fs := tokens.NewFileSet()
	f := fileFor(fs, text)
	start := 0
	l := len(text)
	var errs []errors.Error
	for start < l {
		index := strings.Index(text[start:l], pattern)
		if index >= 0 {
			errs = append(errs, errors.New(location.NewLocation(f.Pos(index+start), f.Pos(index+start+len(pattern))), message))
			start += index + len(pattern)
		} else {
			break
		}
	}
	return fs, p, errs
}

func TestDiagnostics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Diagnositics Suite")
}
