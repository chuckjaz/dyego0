package binder

import (
	"dyego0/assert"
	"dyego0/ast"
	"dyego0/symbols"
	"dyego0/types"
)

type buildVisitor struct {
	container           types.TypeSymbol
	scope               symbols.Scope
	members             []types.Member
	membersScopeBuilder symbols.ScopeBuilder
	typeScopeBuilder    symbols.ScopeBuilder
	signatures          []types.Signature
	builders            map[symbols.Symbol]symbols.ScopeBuilder
	openTypeSymbols     map[types.TypeSymbol]ast.Element
	openElements        map[ast.Element]types.TypeSymbol
	context             *BindingContext
}

func newBuilderVisitor(
	container types.TypeSymbol,
	scope symbols.Scope,
	context *BindingContext,
	builders map[symbols.Symbol]symbols.ScopeBuilder,
	typeScopeBuilder symbols.ScopeBuilder,
) *buildVisitor {
	return &buildVisitor{
		container:           container,
		scope:               scope,
		context:             context,
		builders:            builders,
		membersScopeBuilder: symbols.NewBuilder(),
		openTypeSymbols:     make(map[types.TypeSymbol]ast.Element),
		openElements:        make(map[ast.Element]types.TypeSymbol),
		typeScopeBuilder:    typeScopeBuilder,
	}
}

func (v *buildVisitor) findTypeInType(element ast.Element, typeSym types.TypeSymbol) types.TypeSymbol {
	if types.IsError(typeSym) {
		return typeSym
	}
	t := typeSym.Type()
	if t == nil {
		// Type is not built yet, use the builder instead
		b, ok := v.builders[typeSym]
		assert.Assert(ok, "Unbuilt type not found in builders")
		return v.findTypeIn(element, b)
	}
	return v.findTypeIn(element, t.TypeScope())
}

func (v *buildVisitor) findTypeIn(element ast.Element, scope symbols.Scope) types.TypeSymbol {
	switch n := element.(type) {
	case ast.Name:
		sym, ok := scope.Find(n.Text())
		if !ok {
			v.context.Error(n, "Undefined symbol %s", n.Text())
			return types.NewErrorType()
		}
		typeSym, ok := sym.(types.TypeSymbol)
		if !ok {
			v.context.Error(n, "Expected %s to be a type symbol", n.Text())
			return types.NewErrorType()
		}
		return typeSym
	case ast.Selection:
		container := v.findTypeIn(n.Target(), scope)
		if types.IsError(container) {
			return container
		}
		return v.findTypeInType(n.Member(), container)
	case ast.SequenceType:
		elements := v.findTypeIn(n.Elements(), scope)
		return types.MakeArray(elements)
	case ast.ReferenceType:
		referant := v.findTypeIn(n.Referent(), scope)
		return types.MakeReference(referant)
	}
	assert.Fail("Unhandled element type %#v", element)
	return nil
}

func (v *buildVisitor) findType(element ast.Element) types.TypeSymbol {
	return v.findTypeIn(element, v.scope)
}

func (v *buildVisitor) openTypeFor(element ast.Element) types.TypeSymbol {
	sym, ok := v.openElements[element]
	assert.Assert(!ok, "Duplicate use of element %s", element)
	sym = types.NewTypeSymbol("", nil)
	v.openElements[element] = sym
	v.openTypeSymbols[sym] = element
	return sym
}

func (v *buildVisitor) enterMember(element ast.Element, member types.Member) {
	_, ok := v.membersScopeBuilder.Enter(member)
	if ok {
		v.members = append(v.members, member)
	} else {
		v.context.Error(element, "Duplicate member")
	}
}

func (v *buildVisitor) enterTypeMember(element ast.Element, member types.TypeMember) {
	_, ok := v.typeScopeBuilder.Enter(member)
	if !ok {
		v.context.Error(element, "Duplicate member")
	}
}

func (v *buildVisitor) targetAndContextOf(element ast.Element) (types.TypeSymbol, []types.TypeSymbol) {
	var targetName ast.Name
	var contextNames []ast.Name
	switch n := element.(type) {
	case ast.Name:
		targetName = n
	case ast.Selection:
		current := n.Target()
		targetName = n.Member()
		for true {
			switch m := current.(type) {
			case ast.Name:
				contextNames = append([]ast.Name{m}, contextNames...)
			case ast.Selection:
				contextNames = append([]ast.Name{m.Member()}, contextNames...)
				current = m.Target()
				continue
			default:
				assert.Fail("Unknown node in targetAndContext: %s", current)
			}
			break
		}
	}
	target := v.findType(targetName)
	var context []types.TypeSymbol
	for _, contextName := range contextNames {
		context = append(context, v.findType(contextName))
	}
	return target, context
}

func (v *buildVisitor) Visit(element ast.Element) bool {
	for {
		switch n := element.(type) {
		case ast.Sequence:
			v.Visit(n.Left())
			element = n.Right() // Simulated tail call
			continue
		case ast.Storage:
			var ft types.TypeSymbol
			if n.Type() == nil {
				ft = types.NewTypeSymbol("", nil)
			} else {
				ft = v.findType(n.Type())
			}
			f := types.NewField(n.Name().Text(), ft, n.Mutable())
			v.enterMember(element, f)
		case ast.Definition:
			if isTypeDeclaration(n) {
				name, ok := n.Name().(ast.Name)
				assert.Assert(ok, "Expected types to be a single identifier")
				sym, ok := v.typeScopeBuilder.Find(name.Text())
				assert.Assert(ok, "Type symbol not found")
				typeSym, ok := sym.(types.TypeSymbol)
				assert.Assert(ok, "Expected a type symbol %#v", typeSym)
				builder, ok := v.builders[typeSym]
				assert.Assert(ok, "Build missing")
				nestedScope := symbols.Merge(v.scope, v.typeScopeBuilder)
				nested := newBuilderVisitor(typeSym, nestedScope, v.context, v.builders, builder)
				ast.WalkChildren(n.Value(), nested)
				nested.Done(typeSym, types.Record, v.container)
			} else {
				var typeSym types.TypeSymbol
				if n.Type() != nil {
					typeSym = v.findType(n.Type())
				} else {
					typeSym = v.openTypeFor(n.Value())
				}
				v.enterTypeMember(n, types.NewTypeMember(n.Name().Text(), typeSym))
			}
		}
		break
	}

	return false
}

func (v *buildVisitor) Done(typeSym types.TypeSymbol, kind types.TypeKind, container types.TypeSymbol) {
	types.NewType(
		typeSym,
		kind,
		v.members,
		v.membersScopeBuilder.Build(),
		v.typeScopeBuilder.Build(),
		v.signatures,
		container,
	)
}

// Build builds the types in the given module
func (c *BindingContext) Build(moduleSymbol types.TypeSymbol, element ast.Element) {
	v := newBuilderVisitor(moduleSymbol, c.Scope, c, c.Builders, c.Scope)
	v.Visit(element)
	v.Done(moduleSymbol, types.Module, nil)
}
