package binder

import (
	"dyego0/ast"
	"dyego0/errors"
	"dyego0/symbols"
	"dyego0/types"
)

type enterVisitor struct {
	builder *typeBuilder
	c       *BindingContext
	errors  []errors.Error
}

func newEnterVisitor(builder *typeBuilder, c *BindingContext) *enterVisitor {
	return &enterVisitor{builder: builder, c: c}
}

func (v *enterVisitor) enterSymbol(symbol symbols.Symbol, node ast.Element) {
	ok := v.builder.AddTypeSymbol(symbol)
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
		if isTypeLiteral(n.Value()) {
			v.handleTypeLiteral(n.Name(), n.Value().(ast.TypeLiteral))
		} else {
			typeConstant := types.NewTypeConstant(n.Name().Text(), nil, nil)
			v.enterSymbol(typeConstant, n.Name())
		}
	case ast.Spread:

	case ast.TypeLiteral:
		ast.WalkChildren(n, v)
	}
	return true
}

func (v *enterVisitor) handleTypeLiteral(name ast.Name, literal ast.TypeLiteral) {
	symbol := types.NewTypeSymbol(name.Text(), nil)
	v.enterSymbol(symbol, name)
	nestedBuilder := newTypeBuilder(symbol)
	v.builder.RecordNestedTypeBuilder(symbol, nestedBuilder)
	v.c.Enter(literal, nestedBuilder)
}

func isTypeLiteral(value ast.Element) bool {
	_, ok := value.(ast.TypeLiteral)
	return ok
}

// Enter symbols for the let definition in element into builder
func (c *BindingContext) Enter(element ast.Element, builder *typeBuilder) {
	v := newEnterVisitor(builder, c)
	v.Visit(element)
	c.Errors = append(c.Errors, v.errors...)
}
