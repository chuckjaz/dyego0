package types_test

import (
	"testing"

	"dyego0/symbols"
	"dyego0/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("types", func() {
	It("should be able to create a type", func() {
		s := types.NewTypeSymbol("A", nil)
		t := types.NewType(s, nil, nil, nil, nil)
		Expect(t.Symbol()).To(Equal(s))
		Expect(s.Type()).To(Equal(t))
		Expect(t.DisplayName()).To(Equal("A"))
		Expect(t.Members()).To(BeNil())
		Expect(t.MembersScope()).To(BeNil())
		Expect(t.Signatures()).To(BeNil())
		Expect(t.TypeScope()).To(BeNil())
		Expect(s.Canonical()).To(Equal(s))
		Expect(s.IsType()).To(BeTrue())
	})
	It("should be able the create a field", func() {
		f := types.NewField("a", nil, true)
		Expect(f.Name()).To(Equal("a"))
		Expect(f.Type()).To(BeNil())
		Expect(f.IsField()).To(BeTrue())
		Expect(f.IsMember()).To(BeTrue())
		Expect(f.Mutable()).To(BeTrue())
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
	It("can produce a display name for anonymous types", func() {
		a := types.NewType(types.NewTypeSymbol("A", nil), nil, nil, nil, nil)
		field := types.NewField("a", a, false)
		membersScopeBuilder := symbols.NewBuilder()
		_, ok := membersScopeBuilder.Enter(field)
		Expect(ok).To(BeTrue())
		membersScope := membersScopeBuilder.Build()
		t := types.NewType(
			types.NewTypeSymbol("", nil),
			[]types.Member{types.NewField("a", a, false)},
			membersScope,
			[]types.Signature{
				types.NewSignature(
					a,
					[]types.Parameter{types.NewParameter("a", a)},
					a,
				)},
			nil,
		)
		Expect(t.DisplayName()).To(Equal("<a: A, {A.a: A -> A}>"))
	})
})

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Symbols Suite")
}
