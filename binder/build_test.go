package binder_test

import (
	"dyego0/ast"
	"dyego0/binder"
	"dyego0/symbols"
	"dyego0/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build", func() {
	p := func(text string) ast.Element {
		return parse("let Int = <>\n" + text)
	}
	findTypeMember := func(typ types.Type, name string) symbols.Symbol {
		resultSym, ok := typ.TypeScope().Find(name)
		Expect(ok).To(BeTrue())
		return resultSym
	}
	findType := func(typ types.Type, name string) types.Type {
		resultSym := findTypeMember(typ, name)
		typeSym, ok := resultSym.(types.TypeSymbol)
		Expect(ok).To(BeTrue())
		return typeSym.Type()
	}
	findMember := func(typ types.Type, name string) types.Member {
		result, ok := typ.MemberScope().Find(name)
		Expect(ok).To(BeTrue())
		resultSym, ok := result.(types.Member)
		Expect(ok).To(BeTrue())
		return resultSym
	}
	findExtension := func(typ types.Type, name string) types.TypeExtension {
		result, ok := typ.Extensions().Find(name)
		Expect(ok).To(BeTrue())
		resultSym, ok := result.(types.TypeExtensions)
		Expect(ok).To(BeTrue())
		return resultSym.Extensions()[0]
	}
	m := func(text string) types.Type {
		context := binder.NewContext()
		element := p(text)
		module := types.NewTypeSymbol("moule", nil)
		context.Enter(element)
		context.Build(module, element)
		Expect(context.Errors).To(BeNil())
		return module.Type()
	}
	It("can build a single type", func() {
		module := m("let a = < a: Int >")
		at := findType(module, "a")
		Expect(at).To(Not(BeNil()))
		af := findMember(at, "a")
		Expect(af).To(Not(BeNil()))
	})
	It("can build module literal", func() {
		modules := m("let a = 1")
		am := findTypeMember(modules, "a")
		Expect(am).To(Not(BeNil()))
	})
	It("can build a var", func() {
		modules := m("var a = 1")
		am := findMember(modules, "a")
		Expect(am).To(Not(BeNil()))
	})
	It("can extend a type", func() {
		modules := m("let Int.size = 20")
		e := findExtension(modules, "size")
		Expect(e).To(Not(BeNil()))
		Expect(e.Target().Name()).To(Equal("Int"))
	})
	It("can extend with a context", func() {
		modules := m("let Int.Int.size = 20")
		e := findExtension(modules, "size")
		Expect(e).To(Not(BeNil()))
		Expect(e.Context()[0].Name()).To(Equal("Int"))
	})
	It("can find a nested type", func() {
		modules := m("let A = < B = <> >, var b: A.B")
		mb := findMember(modules, "b")
		Expect(mb).To(Not(BeNil()))
	})
	It("can create a reference type", func() {
		modules := m("var a: *Int")
		ma := findMember(modules, "a")
		Expect(ma).To(Not(BeNil()))
	})
	It("can create a sequence reference", func() {
		modules := m("var a: Int[]")
		ma := findMember(modules, "a")
		Expect(ma).To(Not(BeNil()))
	})
})
