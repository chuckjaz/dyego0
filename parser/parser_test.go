package parser_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego0/ast"
	"dyego0/parser"
)

var _ = Describe("parser", func() {
	Describe("literals", func() {
		It("can parse an int", func() {
			l, ok := parse("123").(ast.LiteralInt)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(123))
		})
	})
})


func parse(text string) ast.Element {
	s := parser.NewScanner(append([]byte(text), 0), 0)
	p := parser.NewParser(s)
	r := p.Parse()
	Expect(p.Errors()).To(BeNil())
	return r
}

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}
