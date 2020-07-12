package ast_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "go/token"
    "dyego0/ast"
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
        It("Selection", func() {
            l := b.Selection(nil, nil)
            Expect(l.Target()).To(BeNil())
            Expect(l.Member()).To(BeNil())
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
            l := b.Lambda(nil, nil)
            Expect(l.Parameters()).To(BeNil())
            Expect(l.Body()).To(BeNil())
        })
        It("Parameter", func() {
            l := b.Parameter(b.Name("name"), nil, nil)
            Expect(l.Name().Text()).To(Equal("name"))
            Expect(l.Type()).To(BeNil())
            Expect(l.Default()).To(BeNil())
            Expect(l.IsParameter()).To(Equal(true))
        })
        It("VarDefinition", func() {
            l := b.VarDefinition(b.Name("name"), nil, false)
            Expect(l.Name().Text()).To(Equal("name"))
            Expect(l.Type()).To(BeNil())
            Expect(l.Mutable()).To(Equal(false))
            Expect(l.IsField()).To(Equal(true))
        })
        It("Error", func() {
            l := b.Error("message")
            Expect(l.Message()).To(Equal("message"))
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
            n  = b.Name("name")
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

