package ast_test

import "go/token"
import . "github.com/onsi/ginkgo"
import . "github.com/onsi/gomega"
import "dyego0/ast"

var _ = Describe("location", func() {
	s := token.Pos(1)
	e := token.Pos(100)
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
