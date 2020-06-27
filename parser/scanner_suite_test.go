package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego/parser"
	"dyego/tokens"
)

func scanBytes(src []byte, expected ...tokens.Token) int {
	scanner := parser.NewScanner(src, 0)
	var received tokens.Token
	for _, token := range expected {
		received = scanner.Next()
		Ω(received).Should(Equal(token))
	}
	Ω(scanner.Next()).Should(Equal(tokens.EOF))
	return scanner.Line()
}

func scanString(text string, expected ...tokens.Token) int {
	src := append([]byte(text), 0)
	return scanBytes(src, expected...)
}

var _ = Describe("scanner", func() {
	Describe("when constructing the instance", func() {
		It("should not panic", func() {
			parser.NewScanner([]byte{0}, 0)
		})
	})
	It("should panic when an empty buffer is passed", func() {
		Ω(func() {
			parser.NewScanner([]byte{}, 0)
		}).Should(Panic())
	})
	It("should panic when a non-null-terminated buffer is passed", func() {
		Ω(func() {
			parser.NewScanner([]byte{'a', 'b', 'c'}, 0)
		})
	})
	Describe("when parsing", func() {
		It("should parse 'ident' as an IDENT", func() {
			scanString("ident", tokens.Identifier)
		})
		It("should parse multiple idents as IDENT", func() {
			scanString("  ident   ident2 _ _12", tokens.Identifier, tokens.Identifier, tokens.Identifier, tokens.Identifier)
		})
		It("should parse basic operators correctly", func() {
			scanString("+ - * /", tokens.Add, tokens.Sub, tokens.Mult, tokens.Div)
		})
		It("should count lines correctly", func() {
			lines := scanString(" \n  \r\n  \n ident ", tokens.Identifier)
			Ω(lines).Should(Equal(4))
		})
	})
})
