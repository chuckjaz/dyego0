package binder

import (
	"dyego0/ast"
)

type typeBuilderVisitor struct {
	c *BindingContext
}

func (v *typeBuilderVisitor) Visit(element ast.Element) bool {
	switch n := element.(type) {
	case ast.Sequence:
		v.Visit(n.Left())
		v.Visit(n.Right())
	case ast.LetDefinition:

	}
	return false
}

func (c *BindingContext) Build(element ast.Element) {
}
