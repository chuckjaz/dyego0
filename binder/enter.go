package binder

import (
	"dyego0/ast"
	"dyego0/errors"
	"dyego0/symbols"
	"dyego0/types"
)

type enterVisitor struct {
	builder symbols.ScopeBuilder
	errors  []errors.Error
}

func newEnterVisitor(builder symbols.ScopeBuilder) *enterVisitor {
	return &enterVisitor{builder: builder}
}

func (v *enterVisitor) enterSymbol(symbol symbols.Symbol, node ast.Element) {
	_, ok := v.builder.Enter(symbol)
	if !ok {
		v.errors = append(v.errors, errors.New(node, "Duplicate symbol"))
	}
}

func (v *enterVisitor) Visit(element ast.Element) bool {
	switch n := element.(type) {
	case ast.Sequence:
		v.Visit(n.Left())
		v.Visit(n.Right())
	case ast.LetDefinition:
		if isTypeDeclaration(n) {
			name, ok := n.Name().(ast.Name)
			if ok {
				v.enterSymbol(types.NewTypeSymbol(name.Text(), nil), n)
			}
		}
	}
	return true
}

func isTypeDeclaration(declaration ast.LetDefinition) bool {
	_, ok := declaration.Value().(ast.TypeLiteral)
	return ok
}

// Enter enters the types declared a the root of emement into the scope
func (c *BindingContext) Enter(element ast.Element, builder symbols.ScopeBuilder) {
	v := newEnterVisitor(builder)
	v.Visit(element)
	c.Errors = append(c.Errors, v.errors...)
}
