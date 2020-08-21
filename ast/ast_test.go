package ast_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego0/ast"
	"go/token"
)

var _ = Describe("ast", func() {
	Describe("construct", func() {
		b := ast.NewBuilder(ast.NewLocation(0, 1))
		b.PushContext()
		It("Name", func() {
			n := b.Name("text")
			Expect(n.Text()).To(Equal("text"))
		})
		It("LiteralRune", func() {
			n := b.LiteralRune('a')
			Expect(n.Value()).To(Equal('a'))
		})
		It("LiteralByte", func() {
			n := b.LiteralByte(byte(42))
			Expect(n.Value()).To(Equal(byte(42)))
		})
		It("LiteralInt", func() {
			n := b.LiteralInt(1)
			Expect(n.Value()).To(Equal(1))
		})
		It("LiteralUInt", func() {
			n := b.LiteralUInt(uint(42))
			Expect(n.Value()).To(Equal(uint(42)))
		})
		It("LiteralLong", func() {
			n := b.LiteralLong(int64(42))
			Expect(n.Value()).To(Equal(int64(42)))
		})
		It("LiteralULong", func() {
			n := b.LiteralULong(uint64(42))
			Expect(n.Value()).To(Equal(uint64(42)))
		})
		It("LiteralDouble", func() {
			l := b.LiteralDouble(1.2)
			Expect(l.Value()).To(BeNumerically("~", 1.2))
		})
		It("LitearlFloat", func() {
			l := b.LiteralFloat(1.2)
			Expect(l.Value()).To(BeNumerically("~", float32(1.2)))
		})
		It("LiteralBoolean", func() {
			l := b.LiteralBoolean(true)
			Expect(l.Value()).To(Equal(true))
		})
		It("LiteralString", func() {
			l := b.LiteralString("value")
			Expect(l.Value()).To(Equal("value"))
		})
		It("LiteralNull", func() {
			l := b.LiteralNull()
			Expect(l.IsNull()).To(Equal(true))
		})
		It("Break", func() {
			n := b.Break(nil)
			Expect(n.Label()).To(BeNil())
			Expect(n.IsBreak()).To(BeTrue())
		})
		It("Continue", func() {
			n := b.Continue(nil)
			Expect(n.Label()).To(BeNil())
			Expect(n.IsContinue()).To(BeTrue())
		})
		It("Selection", func() {
			l := b.Selection(nil, nil)
			Expect(l.Target()).To(BeNil())
			Expect(l.Member()).To(BeNil())
		})
		It("Sequence", func() {
			s := b.Sequence(nil, nil)
			Expect(s.Left()).To(BeNil())
			Expect(s.Right()).To(BeNil())
			Expect(s.IsSequence()).To(Equal(true))
		})
		It("Spread", func() {
			s := b.Spread(nil)
			Expect(s.Target()).To(BeNil())
			Expect(s.IsSpread()).To(Equal(true))
		})
		It("Call", func() {
			l := b.Call(nil, nil)
			Expect(l.Target()).To(BeNil())
			Expect(l.Arguments()).To(BeNil())
		})
		It("NamedArgument", func() {
			l := b.NamedArgument(nil, nil)
			Expect(l.Name()).To(BeNil())
			Expect(l.Value()).To(BeNil())
			Expect(l.IsNamedArgument()).To(Equal(true))
		})
		It("ObjectInitializer", func() {
			l := b.ObjectInitializer(false, nil)
			Expect(l.Mutable()).To(Equal(false))
			Expect(l.Members()).To(BeNil())
			Expect(l.IsObject()).To(Equal(true))
		})
		It("ArrayInitializer", func() {
			l := b.ArrayInitializer(false, nil)
			Expect(l.Mutable()).To(Equal(false))
			Expect(l.Elements()).To(BeNil())
			Expect(l.IsArray()).To(Equal(true))
		})
		It("NameMemberInitializer", func() {
			l := b.NamedMemberInitializer(b.Name("name"), nil, nil)
			Expect(l.Name().Text()).To(Equal("name"))
			Expect(l.Type()).To(BeNil())
			Expect(l.Value()).To(BeNil())
		})
		It("SplatMemberInitializer", func() {
			l := b.SplatMemberInitializer(nil)
			Expect(l.Type()).To(BeNil())
			Expect(l.IsSplat()).To(Equal(true))
		})
		It("Lambda", func() {
			l := b.Lambda(nil, nil, nil)
			Expect(l.TypeParameters()).To(BeNil())
			Expect(l.Parameters()).To(BeNil())
			Expect(l.Body()).To(BeNil())
		})
		It("Loop", func() {
			l := b.Loop(nil, nil)
			Expect(l.Label()).To(BeNil())
			Expect(l.Body()).To(BeNil())
			Expect(l.IsLoop()).To(BeTrue())
		})
		It("Parameter", func() {
			l := b.Parameter(b.Name("name"), nil, nil)
			Expect(l.Name().Text()).To(Equal("name"))
			Expect(l.Type()).To(BeNil())
			Expect(l.Default()).To(BeNil())
			Expect(l.IsParameter()).To(Equal(true))
		})
		It("Return", func() {
			n := b.Return(nil)
			Expect(n.Value()).To(BeNil())
			Expect(n.IsReturn()).To(BeTrue())
		})
		It("TypeParameters", func() {
			l := b.TypeParameters(nil, nil)
			Expect(l.Parameters()).To(BeNil())
			Expect(l.Wheres()).To(BeNil())
		})
		It("TypeParameter", func() {
			l := b.TypeParameter(nil, nil)
			Expect(l.Name()).To(BeNil())
			Expect(l.Constraint()).To(BeNil())
			Expect(l.IsTypeParameter()).To(Equal(true))
		})
		It("When", func() {
			n := b.When(nil, nil)
			Expect(n.Target()).To(BeNil())
			Expect(n.Clauses()).To(BeNil())
		})
		It("WhenElseClause", func() {
			n := b.WhenElseClause(nil)
			Expect(n.Body()).To(BeNil())
			Expect(n.IsElse()).To(BeTrue())
		})
		It("WhenValueClause", func() {
			n := b.WhenValueClause(nil, nil)
			Expect(n.Value()).To(BeNil())
			Expect(n.Body()).To(BeNil())
			Expect(n.IsWhenValueClause()).To(BeTrue())
		})
		It("Where", func() {
			l := b.Where(nil, nil)
			Expect(l.Left()).To(BeNil())
			Expect(l.Right()).To(BeNil())
			Expect(l.IsWhere()).To(Equal(true))
		})
		It("VarDefinition", func() {
			l := b.VarDefinition(b.Name("name"), nil, false)
			Expect(l.Name().Text()).To(Equal("name"))
			Expect(l.Type()).To(BeNil())
			Expect(l.Mutable()).To(Equal(false))
			Expect(l.IsField()).To(Equal(true))
		})
		It("LetDefinition", func() {
			l := b.LetDefinition(nil, nil)
			Expect(l.Name()).To(BeNil())
			Expect(l.Value()).To(BeNil())
			Expect(l.IsLetDefinition()).To(Equal(true))
		})
		It("VocabularyLiteral", func() {
			l := b.VocabularyLiteral(nil)
			Expect(l.Members()).To(BeNil())
			Expect(l.IsVocabularyLiteral()).To(Equal(true))
		})
		It("VocabularyOperatorDeclaration", func() {
			l := b.VocabularyOperatorDeclaration(nil, ast.Infix, nil, ast.Left)
			Expect(l.Names()).To(BeNil())
			Expect(l.Placement()).To(Equal(ast.Infix))
			Expect(l.Precedence()).To(BeNil())
			Expect(l.Associativity()).To(Equal(ast.Left))
		})
		It("VocabularyOperatorPrecedence", func() {
			l := b.VocabularyOperatorPrecedence(nil, ast.Infix, ast.Before)
			Expect(l.Name()).To(BeNil())
			Expect(l.Placement()).To(Equal(ast.Infix))
			Expect(l.Relation()).To(Equal(ast.Before))
		})
		It("VocabularyEmbedding", func() {
			l := b.VocabularyEmbedding(nil)
			Expect(l.Name()).To(BeNil())
			Expect(l.IsVocabularyEmbedding()).To(Equal(true))
		})
		It("Error", func() {
			l := b.Error("message")
			Expect(l.Message()).To(Equal("message"))
		})
		It("DirectError", func() {
			l := b.DirectError(token.Pos(1), token.Pos(2), "message")
			Expect(l.Message()).To(Equal("message"))
			Expect(l.Start()).To(Equal(token.Pos(1)))
			Expect(l.End()).To(Equal(token.Pos(2)))
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
			Expect(n.Start()).To(Equal(token.Pos(10)))
			Expect(n.End()).To(Equal(token.Pos(15)))
			Expect(n.Length()).To(Equal(5))
			b.PopContext()
			n = b.Name("name")
			Expect(n.Start()).To(Equal(token.Pos(0)))
			Expect(n.End()).To(Equal(token.Pos(15)))
			l.start = 20
			l.end = 30
			b.PushContext()
			n = b.Name("name")
			Expect(n.Start()).To(Equal(token.Pos(20)))
			Expect(n.End()).To(Equal(token.Pos(30)))
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
			Expect(bn.Start()).To(Equal(token.Pos(100)))
			Expect(cn.Start()).To(Equal(token.Pos(0)))
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
	start, end token.Pos
}

func (b BuilderContext) Start() token.Pos {
	return b.start
}

func (b BuilderContext) End() token.Pos {
	return b.end
}

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}
