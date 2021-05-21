package binder_test

import (
	"dyego0/ast"
	"dyego0/binder"
	"dyego0/diagnostics"
	"dyego0/parser"
	"dyego0/scanner"
	"dyego0/tokens"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("enter", func() {
	It("should be able to enter a single type", func() {
		context := binder.NewContext()
		element := parse("let a = < a: Int >")
		context.Enter(element)
		Expect(len(context.Errors)).To(Equal(0))
		scope := context.Scope.Build()
		_, ok := scope.Find("a")
		Expect(ok).To(BeTrue())
	})
	It("should be able to enter multiple types", func() {
		context := binder.NewContext()
		element := parse("let a = < a: Int >, let b = < b: Int >")
		context.Enter(element)
		Expect(context.Errors).To(BeNil())
		scope := context.Scope.Build()
		_, ok := scope.Find("a")
		Expect(ok).To(BeTrue())
		_, ok = scope.Find("b")
		Expect(ok).To(BeTrue())
	})
	It("should detect a duplicate symbol", func() {
		context := binder.NewContext()
		element := parse("let a = < a: Int >, let a = < b: Int >")
		context.Enter(element)
		Expect(len(context.Errors)).To(Equal(1))
		Expect(context.Errors[0].Error()).To(Equal("Duplicate symbol"))
	})
	It("should be able to enter nested types", func() {
		context := binder.NewContext()
		element := parse("let a = < a: Int, b = < a: Int > >")
		context.Enter(element)
		aSym, ok := context.Scope.Find("a")
		Expect(ok).To(BeTrue())
		bBuilder, ok := context.Builders[aSym]
		Expect(ok).To(BeTrue())
		_, ok = bBuilder.Find("b")
		Expect(ok).To(BeTrue())
	})
})

func recordLines(fb tokens.FileBuilder, text string) {
	for o, ch := range text {
		if ch == '\n' {
			fb.AddLine(o)
		}
	}
}

type source struct {
	text string
}

func (s source) Source(filename string) diagnostics.Source {
	return s
}

func (s source) Text(start, end int) string {
	return s.text[start:end]
}

func scan(text string, fb tokens.FileBuilder) *scanner.Scanner {
	return scanner.NewScanner(append([]byte(text), 0), 0, fb)
}

func parseNamed(text, filename string) ast.Element {
	fs := tokens.NewFileSet()
	fb := fs.BuildFile(filename, len(text))
	p := parser.NewParser(scan(text, fb), nil)
	r := p.Parse()
	recordLines(fb, text)
	fb.Build()
	errs := p.Errors()
	if errs != nil {
		sf := source{text: text}
		print(diagnostics.Format(errs, fs, sf))
	}
	Expect(p.Errors()).To(BeNil())
	return r
}

func parse(text string) ast.Element {
	return parseNamed(text, "test")
}
