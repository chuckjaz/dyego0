package ast_test

import (
	"dyego0/ast"
	"dyego0/tokens"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("location", func() {
	s := tokens.Pos(1)
	e := tokens.Pos(100)
	l := ast.NewLocation(s, e)

	It("should access start", func() {
		Expect(l.Start()).To(Equal(s))
	})

	It("should access end", func() {
		Expect(l.End()).To(Equal(e))
	})

	It("should access length", func() {
		Expect(l.Length()).To(Equal(99))
	})
})
