package ast_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego0/ast"
	"dyego0/location"
)

func s(a interface{}) string {
	return fmt.Sprintf("%s", a)
}

var _ = Describe("ast", func() {
	Describe("construct", func() {
		b := ast.NewBuilder(location.NewLocation(0, 1))
		b.PushContext()
		It("Name", func() {
			n := b.Name("text")
			Expect(n.Text()).To(Equal("text"))
			Expect(s(n)).To(Equal("Name(Location(0-1), text)"))
		})
		It("Literal", func() {
			n := b.Literal('a')
			Expect(n.Value()).To(Equal('a'))
			Expect(s(n)).To(Equal("Literal(Location(0-1), 'a')"))
		})
		It("Break", func() {
			n := b.Break(nil)
			Expect(n.Label()).To(BeNil())
			Expect(n.IsBreak()).To(BeTrue())
			Expect(s(n)).To(Equal("Break(Location(0-1))"))
		})
		It("Continue", func() {
			n := b.Continue(nil)
			Expect(n.Label()).To(BeNil())
			Expect(n.IsContinue()).To(BeTrue())
			Expect(s(n)).To(Equal("Continue(Location(0-1))"))
		})
		It("Selection", func() {
			l := b.Selection(nil, nil)
			Expect(l.Target()).To(BeNil())
			Expect(l.Member()).To(BeNil())
			Expect(s(l)).To(Equal("Selection(Location(0-1), target: nil, member: nil)"))
		})
		It("Sequence", func() {
			n := b.Sequence(nil, nil)
			Expect(n.Left()).To(BeNil())
			Expect(n.Right()).To(BeNil())
			Expect(n.IsSequence()).To(Equal(true))
			Expect(s(n)).To(Equal("Sequence(Location(0-1), left: nil, right: nil)"))
		})
		It("Spread", func() {
			n := b.Spread(nil)
			Expect(n.Target()).To(BeNil())
			Expect(n.IsSpread()).To(Equal(true))
			Expect(s(n)).To(Equal("Spread(Location(0-1), target: nil)"))
		})
		It("Call", func() {
			var args []ast.Element
			args = append(args, b.Literal(1), b.Literal(2))
			l := b.Call(nil, args)
			Expect(l.Target()).To(BeNil())
			Expect(l.Arguments()).To(Equal(args))
			Expect(s(l)).To(Equal("Call(Location(0-1), target: nil, arguments: [Literal(Location(0-1), 1), Literal(Location(0-1), 2)])"))
		})
		It("NamedArgument", func() {
			l := b.NamedArgument(nil, nil)
			Expect(l.Name()).To(BeNil())
			Expect(l.Value()).To(BeNil())
			Expect(l.IsNamedArgument()).To(Equal(true))
			Expect(s(l)).To(Equal("NamedArgument(Location(0-1), name: nil, value: nil)"))
		})
		It("ObjectInitializer", func() {
			l := b.ObjectInitializer(false, nil, nil)
			Expect(l.Mutable()).To(Equal(false))
			Expect(l.Type()).To(BeNil())
			Expect(l.Members()).To(BeNil())
			Expect(l.IsObject()).To(Equal(true))
			Expect(s(l)).To(Equal("ObjectInitializer(Location(0-1), mutable: false, type: nil, members: [])"))
		})
		It("ArrayInitializer", func() {
			l := b.ArrayInitializer(false, nil, nil)
			Expect(l.Mutable()).To(Equal(false))
			Expect(l.Type()).To(BeNil())
			Expect(l.Elements()).To(BeNil())
			Expect(l.IsArray()).To(Equal(true))
			Expect(s(l)).To(Equal("ArrayInitializer(Location(0-1), mutable: false, type: nil, elements: [])"))
		})
		It("NamedMemberInitializer", func() {
			l := b.NamedMemberInitializer(b.Name("name"), nil, nil)
			Expect(l.Name().Text()).To(Equal("name"))
			Expect(l.Type()).To(BeNil())
			Expect(l.Value()).To(BeNil())
			Expect(s(l)).To(Equal("NamedMemberInitializer(Location(0-1), name: Name(Location(0-1), name), type: nil, value: nil)"))
		})
		It("Lambda", func() {
			l := b.Lambda(nil, nil, nil)
			Expect(l.Parameters()).To(BeNil())
			Expect(l.Body()).To(BeNil())
			Expect(l.Result()).To(BeNil())
			Expect(s(l)).To(Equal("Lambda(Location(0-1), parameters: [], body: nil, result: nil)"))
		})
		It("IntrinsicLambda", func() {
			l := b.IntrinsicLambda(nil, nil, nil)
			Expect(l.Parameters()).To(BeNil())
			Expect(l.Body()).To(BeNil())
			Expect(l.Result()).To(BeNil())
			Expect(l.IsIntrinsicLambda()).To(BeTrue())
			Expect(s(l)).To(Equal("IntrinsicLambda(Location(0-1), parameters: [], body: nil, result: nil)"))
		})
		It("Loop", func() {
			l := b.Loop(nil, nil)
			Expect(l.Label()).To(BeNil())
			Expect(l.Body()).To(BeNil())
			Expect(l.IsLoop()).To(BeTrue())
			Expect(s(l)).To(Equal("Loop(Location(0-1), label: nil, body: nil)"))
		})
		It("Parameter", func() {
			l := b.Parameter(b.Name("name"), nil, nil)
			Expect(l.Name().Text()).To(Equal("name"))
			Expect(l.Type()).To(BeNil())
			Expect(l.Default()).To(BeNil())
			Expect(l.IsParameter()).To(Equal(true))
			Expect(s(l)).To(Equal("Parameter(Location(0-1), name: Name(Location(0-1), name), type: nil, default: nil)"))
		})
		It("Return", func() {
			n := b.Return(nil)
			Expect(n.Value()).To(BeNil())
			Expect(n.IsReturn()).To(BeTrue())
			Expect(s(n)).To(Equal("Return(Location(0-1), value: nil)"))
		})
		It("When", func() {
			n := b.When(nil, nil)
			Expect(n.Target()).To(BeNil())
			Expect(n.Clauses()).To(BeNil())
			Expect(s(n)).To(Equal("When(Location(0-1), target: nil, clauses: [])"))
		})
		It("WhenElseClause", func() {
			n := b.WhenElseClause(nil)
			Expect(n.Body()).To(BeNil())
			Expect(n.IsElse()).To(BeTrue())
			Expect(s(n)).To(Equal("WhenElseClause(Location(0-1), body: nil)"))
		})
		It("WhenValueClause", func() {
			n := b.WhenValueClause(nil, nil)
			Expect(n.Value()).To(BeNil())
			Expect(n.Body()).To(BeNil())
			Expect(n.IsWhenValueClause()).To(BeTrue())
			Expect(s(n)).To(Equal("WhenValueClause(Location(0-1), value: nil, body: nil)"))
		})
		It("Storage", func() {
			l := b.Storage(b.Name("name"), nil, nil, false)
			Expect(l.Name().Text()).To(Equal("name"))
			Expect(l.Type()).To(BeNil())
			Expect(l.Value()).To(BeNil())
			Expect(l.Mutable()).To(Equal(false))
			Expect(s(l)).To(Equal("Storage(Location(0-1), name: Name(Location(0-1), name), type: nil, value: nil, mutable: false)"))
		})
		It("Definition", func() {
			l := b.Definition(nil, nil, nil)
			Expect(l.Name()).To(BeNil())
			Expect(l.Value()).To(BeNil())
			Expect(l.IsDefinition()).To(Equal(true))
			Expect(s(l)).To(Equal("Definition(Location(0-1), name: nil, type: nil, value: nil)"))
		})
		It("TypeLiteral", func() {
			t := b.TypeLiteral(nil)
			Expect(t.Members()).To(BeNil())
			Expect(t.IsTypeLiteral()).To(BeTrue())
			Expect(s(t)).To(Equal("TypeLiteral(Location(0-1), members: [])"))
		})
		It("TypeLiteralConstant", func() {
			t := b.TypeLiteralConstant(nil, nil)
			Expect(t.Name()).To(BeNil())
			Expect(t.Value()).To(BeNil())
			Expect(t.IsTypeLiteralConstant()).To(BeTrue())
			Expect(s(t)).To(Equal("TypeLiteralConstant(Location(0-1), name: nil, value: nil)"))
		})
		It("TypeLiteralMember", func() {
			m := b.TypeLiteralMember(nil, nil)
			Expect(m.Name()).To(BeNil())
			Expect(m.Type()).To(BeNil())
			Expect(m.IsTypeLiteralMember()).To(BeTrue())
			Expect(s(m)).To(Equal("TypeLiteralMember(Location(0-1), name: nil, type: nil)"))
		})
		It("CallableTypeMember", func() {
			m := b.CallableTypeMember(nil, nil)
			Expect(m.Parameters()).To(BeNil())
			Expect(m.Result()).To(BeNil())
			Expect(s(m)).To(Equal("CallableTypeMember(Location(0-1), parameters: [], result: nil)"))
		})
		It("SequenceType", func() {
			n := b.SequenceType(nil)
			Expect(n.Elements()).To(BeNil())
			Expect(n.IsSequenceType()).To(BeTrue())
			Expect(s(n)).To(Equal("SequenceType(Location(0-1), elements: nil)"))
		})
		It("OptionalType", func() {
			n := b.OptionalType(nil)
			Expect(n.Target()).To(BeNil())
			Expect(n.IsOptionalType()).To(BeTrue())
			Expect(s(n)).To(Equal("OptionalType(Location(0-1), target: nil"))
		})
		It("ReferenceType", func() {
			n := b.ReferenceType(nil)
			Expect(n.Referent()).To(BeNil())
			Expect(s(n)).To(Equal("ReferenceType(Location(0-1), referent: nil)"))
		})
		It("VocabularyLiteral", func() {
			l := b.VocabularyLiteral(nil)
			Expect(l.Members()).To(BeNil())
			Expect(l.IsVocabularyLiteral()).To(Equal(true))
			Expect(s(l)).To(Equal("VocabularyLiteral(Location(0-1), members: [])"))
		})
		It("VocabularyOperatorDeclaration", func() {
			var names []ast.Name
			names = append(names, b.Name("a"), b.Name("b"))
			l := b.VocabularyOperatorDeclaration(names, ast.Infix, nil, ast.Left)
			Expect(l.Names()).To(Equal(names))
			Expect(l.Placement()).To(Equal(ast.Infix))
			Expect(l.Precedence()).To(BeNil())
			Expect(l.Associativity()).To(Equal(ast.Left))
			Expect(s(l)).To(
				Equal("VocabularyOperatorDeclaration(Location(0-1), names: [Name(Location(0-1), a), Name(Location(0-1), b)], placement: infix, precedence: nil, associativity: left)"),
			)
		})
		It("VocabularyOperatorPrecedence", func() {
			l := b.VocabularyOperatorPrecedence(nil, ast.Infix, ast.Before)
			Expect(l.Name()).To(BeNil())
			Expect(l.Placement()).To(Equal(ast.Infix))
			Expect(l.Relation()).To(Equal(ast.Before))
			Expect(s(l)).To(Equal("VocabularyOperatorPrecedence(Location(0-1), name: nil, placement: infix, relation: before)"))
		})
		It("VocabularyEmbedding", func() {
			l := b.VocabularyEmbedding(nil)
			Expect(l.Name()).To(BeNil())
			Expect(l.IsVocabularyEmbedding()).To(Equal(true))
			Expect(s(l)).To(Equal("VocabularyEmbedding(Location(0-1), name: [])"))
		})
		It("Error", func() {
			l := b.Error("message")
			Expect(l.Error()).To(Equal("message"))
		})
	})
	Describe("location", func() {
		It("can push a context", func() {
			l := &BuilderContext{0, 0}
			b := ast.NewBuilder(l)
			b.PushContext()
			l.start = 10
			l.end = 15
			b.PushContext()
			n := b.Name("name")
			Expect(n.Start()).To(Equal(location.Pos(10)))
			Expect(n.End()).To(Equal(location.Pos(15)))
			Expect(n.Length()).To(Equal(5))
			b.PopContext()
			n = b.Name("name")
			Expect(n.Start()).To(Equal(location.Pos(0)))
			Expect(n.End()).To(Equal(location.Pos(15)))
			l.start = 20
			l.end = 30
			b.PushContext()
			n = b.Name("name")
			Expect(n.Start()).To(Equal(location.Pos(20)))
			Expect(n.End()).To(Equal(location.Pos(30)))
			b.PopContext()
			b.PopContext()
		})
	})
	Describe("clone", func() {
		It("can clone a builder", func() {
			bc := &BuilderContext{0, 0}
			b := ast.NewBuilder(bc)
			b.PushContext()
			cc := &BuilderContext{0, 0}
			c := b.Clone(cc)
			bc.start = 100
			bc.end = 101
			b.PushContext()
			bn := b.Name("b")
			cn := c.Name("c")
			Expect(bn.Start()).To(Equal(location.Pos(100)))
			Expect(cn.Start()).To(Equal(location.Pos(0)))
		})
	})
	Describe("String()'s", func() {
		It("OperatorAssociativity", func() {
			Expect(ast.Left.String()).To(Equal("left"))
			Expect(ast.Right.String()).To(Equal("right"))
			Expect(ast.OperatorAssociativity(-1).String()).To(Equal("invalid associativity"))
		})
		It("OperatorPrecedenceReleation", func() {
			Expect(ast.Before.String()).To(Equal("before"))
			Expect(ast.After.String()).To(Equal("after"))
			Expect(ast.OperatorPrecedenceRelation(-1).String()).To(Equal("invalid relation"))
		})
		It("OperatorPlacement", func() {
			Expect(ast.Infix.String()).To(Equal("infix"))
			Expect(ast.Prefix.String()).To(Equal("prefix"))
			Expect(ast.Postfix.String()).To(Equal("postfix"))
			Expect(ast.OperatorPlacement(-1).String()).To(Equal("invalid placement"))
		})
	})
})

type BuilderContext struct {
	start, end location.Pos
}

func (b BuilderContext) Start() location.Pos {
	return b.start
}

func (b BuilderContext) End() location.Pos {
	return b.end
}

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ast Suite")
}
