package binder

import (
	"dyego0/ast"
	"dyego0/diagnostics"
	"dyego0/parser"
	"dyego0/scanner"
	"dyego0/tokens"
	"dyego0/types"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("enter", func() {
	enter := func(text string) *typeBuilder {
		context := NewContext()
		element := parse(text)
		typeSymbol := types.NewTypeSymbol("someType", nil)
		builder := newTypeBuilder(typeSymbol)
		context.Enter(element, builder)
		Expect(len(context.Errors)).To(Equal(0))
		return builder
	}
	It("should be able to enter a single type", func() {
		scope := enter("let a = < a: Int >")
		_, ok := scope.FindTypeSymbol("a")
		Expect(ok).To(BeTrue())
	})
	It("should be able to enter multiple types", func() {
		scope := enter("let a = < a: Int >, let b = < b: Int >")
		_, ok := scope.FindTypeSymbol("a")
		Expect(ok).To(BeTrue())
		_, ok = scope.FindTypeSymbol("b")
		Expect(ok).To(BeTrue())
	})
	It("should detect a duplicate symbol", func() {
		context := NewContext()
		element := parse("let a = < a: Int >, let a = < b: Int >")
		builder := newTypeBuilder(nil)
		context.Enter(element, builder)
		Expect(len(context.Errors)).To(Equal(1))
		Expect(context.Errors[0].Error()).To(Equal("Duplicate symbol"))
	})
	It("should be able to enter sub-types", func() {
		typeBuilder := enter("let a = < a = <a: Int> >")
		topA, ok := typeBuilder.FindTypeSymbol("a")
		Expect(ok).To(BeTrue())
		topABuilder, ok := typeBuilder.FindNestedTypeBuilder(topA)
		Expect(ok).To(BeTrue())
		nestedA, ok := topABuilder.FindTypeSymbol("a")
		Expect(ok).To(BeTrue())
		Expect(nestedA.Name()).To(Equal("a"))
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

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Binding Suite")
}
