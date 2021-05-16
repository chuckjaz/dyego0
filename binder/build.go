package binder

import (
	"dyego0/assert"
	"dyego0/ast"
	"dyego0/errors"
	"dyego0/location"
	"dyego0/symbols"
	"dyego0/types"
)

// typeCache

type typeCache struct {
	c           *BindingContext
	scope       symbols.Scope
	errorType   types.Type
	builtins    types.TypeSymbol
	builtinType map[string]types.TypeSymbol
}

func (c *typeCache) ErrorType() types.TypeSymbol {
	if c.errorType == nil {
		c.errorType = types.NewErrorType()
	}
	return c.errorType.Symbol()
}

func (c *typeCache) error(loc location.Locatable, message string, args ...interface{}) {
	c.c.Errors = append(c.c.Errors, errors.New(loc, message, args...))
}

func (c *typeCache) findType(name string, loc location.Locatable) types.TypeSymbol {
	symbol, ok := c.c.Scope.Find(name)
	var result types.TypeSymbol
	if !ok {
		c.error(loc, "Type '%s' not found", name)
		result = c.ErrorType()
	} else {
		typeSymbol, ok := symbol.(types.TypeSymbol)
		if !ok {
			c.error(loc, "Expected '%s' to be a type", name)
			result = c.ErrorType()
		}
		result = typeSymbol
	}
	return result
}

func (c *typeCache) findBuiltinType(name string, loc location.Locatable) types.TypeSymbol {
	typSym, ok := c.builtinType[name]
	if ok {
		return typSym
	}
	if c.builtins == nil {
		c.builtins = c.findType("$builtins", loc)
	}
	t := c.builtins.Type()
	if t != nil {
		s := t.TypeScope()
		if s != nil {
			sym, ok := s.Find(name)
			if ok {
				typSym, ok = sym.(types.TypeSymbol)
				if ok {
					c.builtinType[name] = typSym
					return typSym
				}
			}
		}
	}

	c.error(loc, "Could not find builtins type %s", name)
	return c.ErrorType()
}

// typeWorker

type typeWorker struct {
	c       *BindingContext
	pending map[symbols.Symbol][]func()
}

func newTypeWorker(c *BindingContext) *typeWorker {
	return &typeWorker{
		c:       c,
		pending: make(map[symbols.Symbol][]func()),
	}
}

func (w *typeWorker) Pend(typed symbols.Symbol, work func()) {
	existing, _ := w.pending[typed]
	w.pending[typed] = append(existing, work)
}

func (w *typeWorker) Resolved(typed symbols.Symbol) {
	work, ok := w.pending[typed]
	if ok {
		delete(w.pending, typed)
		for _, item := range work {
			item()
		}
	}
}

// typeBuilderVisitor

type typeBuilderVisitor struct {
	c           *BindingContext
	typeBuilder *typeBuilder
	typeCache   *typeCache
	typeWorker  *typeWorker
}

func newTypeBuilderVisitor(
	c *BindingContext,
	typeBuilder *typeBuilder,
	typeCache *typeCache,
	typeWorker *typeWorker,
) *typeBuilderVisitor {
	return &typeBuilderVisitor{
		c:           c,
		typeBuilder: typeBuilder,
		typeCache:   typeCache,
		typeWorker:  typeWorker,
	}
}

func (v *typeBuilderVisitor) error(loc location.Locatable, message string, args ...interface{}) {
	v.typeCache.error(loc, message, args...)
}

func (v *typeBuilderVisitor) findType(name string, loc location.Location) types.TypeSymbol {
	return v.typeCache.findType(name, loc)
}

func (v *typeBuilderVisitor) findTypeExpr(element ast.Element) types.TypeSymbol {
	switch n := element.(type) {
	case ast.Selection:
		containerTypeSymbol := v.findTypeExpr(n.Target())
		nestedTypeBuilder, ok := v.typeBuilder.FindNestedTypeBuilder(containerTypeSymbol)
		if !ok {
			v.error(n.Target(), "Expected a type reference")
			return v.typeCache.ErrorType()
		}
		nestedType, ok := nestedTypeBuilder.FindTypeSymbol(n.Member().Text())
		if !ok {
			v.error(n.Target(), "Expected a type reference")
			return v.typeCache.ErrorType()
		}
		nestedTypeSymbol, ok := nestedType.(types.TypeSymbol)
		if !ok {
			v.error(n.Target(), "Expected a type reference")
			return v.typeCache.ErrorType()
		}
		return nestedTypeSymbol
	case ast.Name:
		nestedType, ok := v.typeBuilder.FindTypeSymbol(n.Text())
		if !ok {
			v.error(n, "Expected a type reference")
			return v.typeCache.ErrorType()
		}
		nestedTypeSymbol, ok := nestedType.(types.TypeSymbol)
		if !ok {
			v.error(n, "Expected a type reference")
			return v.typeCache.ErrorType()
		}
		return nestedTypeSymbol
	}
	v.error(element, "Unexpected node", element)
	return v.typeCache.ErrorType()
}

func (v *typeBuilderVisitor) Visit(element ast.Element) bool {
	switch n := element.(type) {
	case ast.Sequence:
		v.Visit(n.Left())
		v.Visit(n.Right())
	case ast.LetDefinition:
		if isTypeLiteral(n.Value()) {
			v.buildTypeFrom(n.Name(), n.Value())
		} else {
			v.buildDefinitionFrom(n.Name(), n.Value())
		}
	case ast.VarDefinition:
		name := n.Name().Text()
		typeSymbol := v.findTypeExpr(n.Type())
		member := types.NewField(name, typeSymbol.Type(), n.Mutable())
		ok := v.typeBuilder.AddMember(member)
		if !ok {
			v.error(n.Name(), "Duplicate member symbol")
		} else {
			v.resolve(member, typeSymbol, n.Name())
		}
	default:
		v.error(element, "Unexpected node", element)
	}
	return false
}

func (v *typeBuilderVisitor) resolve(updatable interface{}, typSym types.TypeSymbol, element ast.Element) {
	updater, ok := updatable.(types.UpdateableType)
	if !ok {
		v.error(element, "Symbol is not updatable")
		return
	}
	typ := typSym.Type()
	if typ == nil {
		v.typeWorker.Pend(typSym, func() {
			updater.UpdateType(typSym.Type())
			sym, ok := updatable.(symbols.Symbol)
			if ok {
				v.typeWorker.Resolved(sym)
			}
		})
	} else {
		updater.UpdateType(typ)
	}
}

func (v *typeBuilderVisitor) inferType(updatable interface{}, element ast.Element) {
	updater, ok := updatable.(types.UpdateableType)
	if !ok {
		v.error(element, "Symbol is not updatable")
		return
	}

}
func (v *typeBuilderVisitor) buildTypeFrom(name ast.Name, element ast.Element) {
	symbol, ok := v.typeBuilder.FindTypeSymbol(name.Text())
	assert.Assert(ok, "Symbol '%s' not entered", name.Text())
	typeSymbol, ok := symbol.(types.TypeSymbol)
	assert.Assert(ok, "Symbol '%s' expected to be a type symbol", name.Text())
	nestedBuilder, ok := v.typeBuilder.FindNestedTypeBuilder(typeSymbol)
	assert.Assert(ok, "Should have been entered by Enter()")
	typeVisitor := newTypeBuilderVisitor(v.c, nestedBuilder, v.typeCache, v.typeWorker)
	typeVisitor.Visit(element)
	typeVisitor.Done()
}

func (v *typeBuilderVisitor) buildDefinitionFromConstExpr(name ast.Name, element ast.Element) {
	symbol, ok := v.typeBuilder.FindTypeSymbol(name.Text())
	assert.Assert(ok, "Symbol '%s' not entered", name.Text())
	constEval := newConstEval(v.typeCache)
	constResult := constEval.Eval(element)
	if constResult.typeCode < tError {
		updatableValue, ok := symbol.(types.UpdateableTypeConstant)
		assert.Assert(ok, "Symbol %s expected to be a type constant", sym)
		updatableValue.UpdateValue(constResult.value)
		v.resolve(symbol, constResult.typ, element)
	} else if constResult.typeCode == tNotConst {
		v.error(element, "Expected constant expression")
	}
}

func (v *typeBuilderVisitor) buildDefinitionFromLambda(name ast.Name, lambda ast.Lambda) {
	thisType := v.typeBuilder.symbol
	var parameters []types.Parameter
	for _, parameter := range lambda.Parameters() {
		parameterName := parameter.Name().Text()
		parameterType := v.findTypeExpr(parameter.Type())
		parameterSym := types.NewParameter(parameterName, nil)
		parameters = append(parameters, parameterSym)
		v.resolve(parameterSym, parameterType, parameter.Type())
	}
	signature := types.NewSignature(nil, parameters, nil)
	v.resolve(signature, thisType, lambda)
	v.inferType(types.ResultTypeUpdater(signature), lambda)
}

func (v *typeBuilderVisitor) buildDefinitionFrom(name ast.Name, element ast.Element) {
	switch n := element.(type) {
	case ast.Lambda:
		v.buildDefinitionFromLambda(name, n)
	default:
		v.buildDefinitionFromConstExpr(name, element)
	}
}

func (v *typeBuilderVisitor) Done() types.Type {
	return v.typeBuilder.Build()
}

func (c *BindingContext) Build(
	element ast.Element,
	moduleSymbol types.TypeSymbol,
	typeBuilder *typeBuilder,
) {
	typeCache := &typeCache{}
	typeWorker := newTypeWorker(c)
	v := newTypeBuilderVisitor(c, typeBuilder, typeCache, typeWorker)
	v.Visit(element)
	v.Done()
}
