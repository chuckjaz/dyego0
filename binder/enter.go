package binder

import (
	"dyego0/ast"
	"dyego0/errors"
	"dyego0/symbols"
	"dyego0/types"
)

type enterVisitor struct {
	scope    symbols.ScopeBuilder
	builders map[symbols.Symbol]symbols.ScopeBuilder
	errors   []errors.Error
}

func newEnterVisitor(scope symbols.ScopeBuilder, builders map[symbols.Symbol]symbols.ScopeBuilder) *enterVisitor {
	return &enterVisitor{scope: scope, builders: builders}
}

func (v *enterVisitor) enterSymbol(symbol symbols.Symbol, node ast.Element) {
	_, ok := v.scope.Enter(symbol)
	if !ok {
		v.errors = append(v.errors, errors.New(node, "Duplicate symbol"))
	}
}

func (v *enterVisitor) Visit(element ast.Element) bool {
	for {
		switch n := element.(type) {
		case ast.Sequence:
			v.Visit(n.Left())
			element = n.Right() // Simulated tail call
			continue
		case ast.Definition:
			typ, ok := n.Value().(ast.TypeLiteral)
			if ok {
				name, ok := n.Name().(ast.Name)
				if ok {
					typSym := types.NewTypeSymbol(name.Text(), nil)
					v.enterSymbol(typSym, n)
					typeScope := symbols.NewBuilder()
					v.builders[typSym] = typeScope
					nestedEnter := newEnterVisitor(typeScope, v.builders)
					for _, member := range typ.Members() {
						nestedEnter.Visit(member)
					}
				}
			}
		}
		break
	}
	return true
}

func isTypeDeclaration(declaration ast.Definition) bool {
	_, ok := declaration.Value().(ast.TypeLiteral)
	return ok
}

// Enter enters the types declared a the root of emement into the scope
func (c *BindingContext) Enter(element ast.Element) {
	v := newEnterVisitor(c.Scope, c.Builders)
	v.Visit(element)
	c.Errors = append(c.Errors, v.errors...)
}
