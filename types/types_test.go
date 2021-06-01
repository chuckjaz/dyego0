package types_test

import (
	"testing"

	"dyego0/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("types", func() {
	It("should be able to create a type", func() {
		s := types.NewTypeSymbol("A", nil)
		t := types.NewType(s, types.Record, nil, nil, nil, nil, nil, nil)
		Expect(t.Symbol()).To(Equal(s))
		Expect(s.Type()).To(Equal(t))
		Expect(t.DisplayName()).To(Equal("A"))
		Expect(t.String()).To(Equal("A"))
		Expect(t.Members()).To(BeNil())
		Expect(t.MemberScope().IsEmpty()).To(BeTrue())
		Expect(t.TypeScope().IsEmpty()).To(BeTrue())
		Expect(t.Signatures()).To(BeNil())
		Expect(s.Canonical()).To(Equal(s))
		Expect(t.Container()).To(BeNil())
		Expect(t.Elements()).To(BeNil())
		Expect(t.Size()).To(Equal(0))
		Expect(t.Referant()).To(BeNil())
		Expect(s.IsType()).To(BeTrue())
	})
	It("should be able the create a field", func() {
		f := types.NewField("a", nil, false)
		Expect(f.Name()).To(Equal("a"))
		Expect(f.Type()).To(BeNil())
		Expect(f.Mutable()).To(BeFalse())
		Expect(f.IsField()).To(BeTrue())
		Expect(f.IsMember()).To(BeTrue())
	})
	It("should be able to create a signature", func() {
		s := types.NewSignature(nil, nil, nil)
		Expect(s.This()).To(BeNil())
		Expect(s.Parameters()).To(BeNil())
		Expect(s.Result()).To(BeNil())
	})
	It("should be able to create a parameter", func() {
		p := types.NewParameter("a", nil)
		Expect(p.Name()).To(Equal("a"))
		Expect(p.Type()).To(BeNil())
		Expect(p.IsParameter()).To(BeTrue())
	})
	It("should be able to create a type member", func() {
		t := types.NewTypeMember("a", nil)
		Expect(t.Name()).To(Equal("a"))
		Expect(t.Type()).To(BeNil())
		Expect(t.IsTypeMember()).To(BeTrue())
	})
	It("can produce a display name for anonymous types", func() {
		a := types.NewTypeSymbol("A", nil)
		types.NewType(a, types.Record, nil, nil, nil, nil, nil, nil)
		t := types.NewType(
			types.NewTypeSymbol("", nil),
			types.Record,
			[]types.Member{types.NewField("a", a, false)},
			nil,
			nil,
			[]types.Signature{types.NewSignature(a, []types.Parameter{types.NewParameter("a", a)}, a)},
			nil,
			nil,
		)
		Expect(t.DisplayName()).To(Equal("<a: A, {A.a: A -> A}>"))
	})
	It("can produce a display name for a contained type", func() {
		c := types.NewTypeSymbol("C", nil)
		types.NewType(c, types.Record, nil, nil, nil, nil, nil, nil)
		n := types.NewTypeSymbol("N", nil)
		t := types.NewType(n, types.Record, nil, nil, nil, nil, nil, c)
		Expect(t.DisplayName()).To(Equal("C.N"))
	})
	It("can produce dislay names for incomplete types", func() {
		e := types.NewTypeSymbol("", nil)
		c := types.NewTypeSymbol("C", nil)
		ne := types.NewTypeSymbol("NE", nil)
		net := types.NewType(ne, types.Record, nil, nil, nil, nil, nil, e)
		Expect(net.DisplayName()).To(Equal("NE"))
		nc := types.NewTypeSymbol("NC", nil)
		nct := types.NewType(nc, types.Record, nil, nil, nil, nil, nil, c)
		Expect(nct.DisplayName()).To(Equal("C.NC"))
	})
	t := func(name string) types.TypeSymbol {
		typeSym := types.NewTypeSymbol(name, nil)
		types.NewType(typeSym, types.Record, nil, nil, nil, nil, nil, nil)
		return typeSym
	}
	It("can make an error types", func() {
		e := types.NewErrorType()
		Expect(types.IsError(e)).To(BeTrue())
	})
	It("can make a sequence", func() {
		e := t("Int")
		a := types.MakeArray(e)
		at := a.Type()
		Expect(at).To(Not(BeNil()))
		Expect(at.Symbol()).To(Equal(a))
		Expect(at.Kind()).To(Equal(types.Array))
		Expect(at.Container()).To(BeNil())
		Expect(at.DisplayName()).To(Equal("Int[]"))
		Expect(at.String()).To(Equal("Int[]"))
		Expect(at.Elements()).To(Equal(e))
		Expect(at.Members()).To(BeNil())
		Expect(at.MemberScope().IsEmpty()).To(BeTrue())
		Expect(at.TypeScope().IsEmpty()).To(BeTrue())
		Expect(at.Size()).To(Equal(-1))
		Expect(at.Referant()).To(BeNil())
		Expect(at.Signatures()).To(BeNil())
	})
	It("can create an array type", func() {
		a := types.NewTypeSymbol("", nil)
		i := t("Int")
		at := types.NewArrayType(a, i, 10)
		Expect(at.DisplayName()).To(Equal("Int[10]"))
	})
	It("can make a referance type", func() {
		i := t("Int")
		r := types.MakeReference(i)
		rt := r.Type()
		Expect(rt).To(Not(BeNil()))
		Expect(rt.Symbol()).To(Equal(r))
		Expect(rt.Kind()).To(Equal(types.Reference))
		Expect(rt.Container()).To(BeNil())
		Expect(rt.DisplayName()).To(Equal("*Int"))
		Expect(rt.String()).To(Equal("*Int"))
		Expect(rt.Elements()).To(BeNil())
		Expect(rt.Members()).To(BeNil())
		Expect(rt.MemberScope().IsEmpty()).To(BeTrue())
		Expect(rt.TypeScope().IsEmpty()).To(BeTrue())
		Expect(rt.Size()).To(Equal(0))
		Expect(rt.Referant()).To(Equal(i))
		Expect(rt.Signatures()).To(BeNil())
	})
})

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Symbols Suite")
}
