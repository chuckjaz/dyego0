package ast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego0/ast"
	"dyego0/location"
)

var _ = Describe("walk", func() {
	b := ast.NewBuilder(location.NewLocation(0, 1))
	b.PushContext()

	n := b.Name("n")
	It("Name", func() {
		expect(n)
	})
	It("Literal", func() {
		expect(b.Literal('a'))
	})
	brk := b.Break(n)
	It("Break", func() {
		expect(brk, n)
	})
	ctn := b.Continue(n)
	It("Continue", func() {
		expect(ctn, n)
	})
	It("Sequence", func() {
		expect(b.Sequence(brk, ctn), brk, n, ctn, n)
	})
	m := b.Name("m")
	It("Selection", func() {
		expect(b.Selection(n, m), n, m)
	})
	It("Spread", func() {
		expect(b.Spread(n), n)
	})
	It("Call", func() {
		expect(b.Call(n, []ast.Element{m, m}), n, m, m)
	})
	one := b.Literal(1)
	It("NamedArgument", func() {
		expect(b.NamedArgument(n, one), n, one)
	})
	It("ObjectInitializer", func() {
		expect(b.ObjectInitializer(true, n, []ast.Element{m, one}), n, m, one)
	})
	It("ArrayInitalizer", func() {
		expect(b.ArrayInitializer(true, n, []ast.Element{m, one}), n, m, one)
	})
	It("NamedMemberInializer", func() {
		expect(b.NamedMemberInitializer(n, m, one), n, m, one)
	})
	param := b.Parameter(n, m, one)
	It("Lambda", func() {
		expect(b.Lambda([]ast.Parameter{param}, one, nil), param, n, m, one, one)
	})
	It("IntrinsicLambda", func() {
		expect(b.IntrinsicLambda([]ast.Parameter{param}, one, m), param, n, m, one, one, m)
	})
	It("Loop", func() {
		expect(b.Loop(n, one), n, one)
	})
	It("Return", func() {
		expect(b.Return(one), one)
	})
	It("When", func() {
		expect(b.When(n, []ast.Element{one}), n, one)
	})
	It("WhenValueClause", func() {
		expect(b.WhenValueClause(n, one), n, one)
	})
	It("WhenElseClause", func() {
		expect(b.WhenElseClause(one), one)
	})
	It("Parameter", func() {
		expect(param, n, m, one)
	})
	It("Definition", func() {
		expect(b.Definition(n, m, one), n, m, one)
	})
	It("Storage", func() {
		expect(b.Storage(n, nil, one, false), n, one)
	})
	It("TypeLiteral", func() {
		expect(b.TypeLiteral([]ast.Element{n}), n)
	})
	It("CallableTypeMember", func() {
		expect(b.CallableTypeMember([]ast.Element{param}, m), param, n, m, one, m)
	})
	It("SequenceType", func() {
		expect(b.SequenceType(n), n)
	})
	It("ReferenceType", func() {
		expect(b.ReferenceType(n), n)
	})
	It("VocabularyLiteral", func() {
		expect(b.VocabularyLiteral([]ast.Element{n}), n)
	})
	precedence := b.VocabularyOperatorPrecedence(n, ast.Infix, ast.Before)
	It("VocabularyOperatorDeclaration", func() {
		expect(b.VocabularyOperatorDeclaration([]ast.Name{n, m}, ast.Infix, precedence, ast.Right), n, m, precedence, n)
	})
	It("VocabularyOperatorPrecedence", func() {
		expect(precedence, n)
	})
	It("VocabularyEmbedding", func() {
		expect(b.VocabularyEmbedding([]ast.Name{n, m}), n, m)
	})
	It("Error", func() {
		expect(b.Error("msg"))
	})
	It("Can stop walking elements early", func() {
		stopsAfter(b.VocabularyLiteral([]ast.Element{n, n, n}), 3)
	})
	It("Can stop walking parameters early", func() {
		stopsAfter(b.Lambda([]ast.Parameter{param, param, param}, one, nil), 3)
	})
	It("Can stop walking names early", func() {
		stopsAfter(b.VocabularyEmbedding([]ast.Name{n, m, n, m}), 3)
	})
})

type testVisitor struct {
	elements []ast.Element
}

func (v *testVisitor) Visit(element ast.Element) bool {
	v.elements = append(v.elements, element)
	return true
}

func expect(element ast.Element, elements ...ast.Element) {
	v := &testVisitor{}
	ast.Walk(element, v)
	expected := append([]ast.Element{element}, elements...)
	Expect(v.elements).To(Equal(expected))
}

type countingVisitor struct {
	count int
}

func (v *countingVisitor) Visit(element ast.Element) bool {
	v.count--
	return v.count > 0
}

func stopsAfter(element ast.Element, n int) {
	v := &countingVisitor{count: n}
	ast.Walk(element, v)
	Expect(v.count).To(Equal(0))
}
