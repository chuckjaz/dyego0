package symbols_test

import (
	"sort"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego0/symbols"
)

type fakeSymbol struct {
	name string
}

func (f fakeSymbol) Name() string {
	return f.name
}

func sym(name string) symbols.Symbol {
	return fakeSymbol{name: name}
}

func scope(names ...string) symbols.Scope {
	b := symbols.NewBuilder()
	for _, name := range names {
		_, ok := b.Enter(sym(name))
		Expect(ok).To(BeTrue())
	}
	return b.Build()
}

func symsFrom(scope symbols.Scope, names ...string) symArray {
	var result symArray
	for _, name := range names {
		s, ok := scope.Find(name)
		Expect(ok).To(BeTrue())
		result = append(result, s)
	}
	return result
}

func expect(scope symbols.Scope, names ...string) {
	for _, name := range names {
		sym, ok := scope.Find(name)
		Expect(ok).To(BeTrue())
		Expect(scope.Contains(name)).To(BeTrue())
		Expect(sym.Name()).To(Equal(name))
	}
}

var _ = Describe("symbols", func() {
	It("can create a builder", func() {
		symbols.NewBuilder()
	})
	It("can enter a symbol", func() {
		b := symbols.NewBuilder()
		_, ok := b.Enter(sym("a"))
		Expect(ok).To(BeTrue())
	})
	It("can build a scope", func() {
		scope("a", "b", "c")
	})
	It("can find a symbol in a scope", func() {
		s := scope("a", "b", "c")
		expect(s, "a", "b", "c")
		_, ok := s.Find("d")
		Expect(ok).To(BeFalse())
	})
	It("can detect a duplicate symbol", func() {
		b := symbols.NewBuilder()
		a := sym("a")
		b.Enter(a)
		previous, ok := b.Enter(sym("a"))
		Expect(ok).To(BeFalse())
		Expect(previous).To(Equal(a))

	})
	It("can reenter a symbol", func() {
		b := symbols.NewBuilder()
		a := sym("a")
		na := sym("a")
		b.Enter(a)
		b.Reenter(na)
		s := b.Build()
		fa, ok := s.Find("a")
		Expect(ok).To(BeTrue())
		Expect(fa).To(Equal(na))
	})
	It("merge symbols", func() {
		s1 := scope("a", "b", "c")
		s2 := scope("d", "e", "f")
		sm := symbols.Merge(s1, s2)
		expect(sm, "a", "b", "c", "d", "e", "f")
	})
	It("can merge to an empty scope", func() {
		s := symbols.Merge()
		_, ok := s.Find("a")
		Expect(ok).To(BeFalse())
		Expect(s.Contains("b")).To(BeFalse())
	})
	It("merging a single scope is a noop", func() {
		s := scope("a", "b", "c")
		sm := symbols.Merge(s)
		Expect(sm).To(Equal(s))
	})
	It("can build one scope from another", func() {
		s1 := scope("a", "b")
		b := symbols.NewBuilderFrom(s1)
		b.Enter(sym("c"))
		s2 := b.Build()
		expect(s2, "a", "b", "c")
	})
	It("can build from a merged scope", func() {
		s1 := symbols.Merge(scope("a"), scope("b"), scope("c"))
		b := symbols.NewBuilderFrom(s1)
		_, ok := b.Enter(sym("a"))
		Expect(ok).To(BeFalse())
		b.Enter(sym("d"))
		expect(b.Build(), "a", "b", "c", "d")
	})
	It("can enumerate a single scope", func() {
		s := scope("a", "b", "c")
		var syms symArray
		s.ForEach(func(symbol symbols.Symbol) bool {
			syms = append(syms, symbol)
			return false
		})
		sort.Sort(syms)
		Expect(syms).To(Equal(symsFrom(s, "a", "b", "c")))
	})
	It("can terminate ForEach early", func() {
		s := scope("a", "b", "c")
		var sym symbols.Symbol
		s.ForEach(func(symbol symbols.Symbol) bool {
			if symbol.Name() == "b" {
				sym = symbol
				return true
			}
			return false
		})
		expected, _ := s.Find("b")
		Expect(sym).To(Equal(expected))
	})
	It("can enumerate a multiScope", func() {
		s := symbols.Merge(scope("a"), scope("b"), scope("c"))
		var syms symArray
		s.ForEach(func(symbol symbols.Symbol) bool {
			syms = append(syms, symbol)
			return false
		})
		Expect(syms).To(Equal(symsFrom(s, "a", "b", "c")))
	})
	It("it ForEach's the correct duplicate", func() {
		s := symbols.Merge(scope("a"), scope("b"), scope("a"), scope("c"))
		var syms symArray
		s.ForEach(func(symbol symbols.Symbol) bool {
			syms = append(syms, symbol)
			return false
		})
		Expect(syms).To(Equal(symsFrom(s, "a", "b", "c")))
	})
	It("can terminate a mult-scope ForEach early", func() {
		s := symbols.Merge(scope("a"), scope("b"), scope("c"))
		var sym symbols.Symbol
		s.ForEach(func(symbol symbols.Symbol) bool {
			if symbol.Name() == "b" {
				sym = symbol
				return true
			}
			return false
		})
		expected, _ := s.Find("b")
		Expect(sym).To(Equal(expected))
	})
})

type symArray []symbols.Symbol

func (a symArray) Len() int {
	return len(a)
}

func (a symArray) Less(i, j int) bool {
	return a[i].Name() < a[j].Name()
}

func (a symArray) Swap(i, j int) {
	tmp := a[i]
	a[i] = a[j]
	a[j] = tmp
}

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Symbols Suite")
}
