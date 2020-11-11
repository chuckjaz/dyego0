package binder

import (
	"dyego0/ast"
	"dyego0/errors"
	"dyego0/symbols"
	"dyego0/types"
)

type enterVisitor struct {
	scope  symbols.ScopeBuilder
	errors []errors.Error
}

func newEnterVisitor(scope symbols.ScopeBuilder) *enterVisitor {
	return &enterVisitor{scope: scope}
}

func (v *enterVisitor) enterSymbol(symbol symbols.Symbol, node ast.Element) {
	_, ok := v.scope.Enter(symbol)
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
func (c *BindingContext) Enter(element ast.Element) {
	v := newEnterVisitor(c.Scope)
	v.Visit(element)
	c.Errors = append(c.Errors, v.errors...)
}
