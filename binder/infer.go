package binder

import (
	"dyego0/assert"
	"dyego0/ast"
	"dyego0/symbols"
	"dyego0/types"
)

type inferencer struct {
	scope        symbols.Scope
	typeWorker   *typeWorker
	typeCache    *typeCache
	elementTypes *elementTypes
	iw           *inferenceWorker
}

func (i *inferencer) Infer(element ast.Element) {
	recordBuiltin := func(name string) {
		typ := i.typeCache.findBuiltinType(name, element)
		i.iw.Record(element, typ)
	}
	switch n := element.(type) {
	case ast.LiteralRune:
		recordBuiltin("Rune")
	case ast.LiteralBoolean:
		recordBuiltin("Boolean")
	case ast.LiteralByte:
		recordBuiltin("Byte")
	case ast.LiteralDouble:
		recordBuiltin("Double")
	case ast.LiteralFloat:
		recordBuiltin("Float")
	case ast.LiteralInt:
		recordBuiltin("Int")
	case ast.LiteralLong:
		recordBuiltin("Long")
	case ast.LiteralULong:
		recordBuiltin("ULong")
	case ast.LiteralNull:
		recordBuiltin("Null")
	case ast.LiteralString:
		recordBuiltin("String")
	case ast.LiteralUInt:
		recordBuiltin("UInt")
	case ast.Selection:
		i.inferTypesOf(func(typs []types.Type) {
			typ := typs[0]
			member, ok := typ.MembersScope().Find(n.Member().Text())
			if !ok {
				i.typeCache.error(n.Member(), "Member not found")
			} else {
				typedMember, ok := member.(types.Member)
				assert.Assert(ok, "Assumed member was a type member")
				i.typeWorker.Pend(typedMember, func() {
					i.iw.Record(element, typedMember.Type().Symbol())
				})
			}
		}, n.Target())
	case ast.Spread:
		recordBuiltin("Unit")
	case ast.Break:
		recordBuiltin("Unit")
	case ast.Call:
		var elements []ast.Element
		elements = append(append(elements, n.Target()), n.Arguments()...)
		i.inferTypesOf(func(typs []types.Type) {
			i.resolveCall(n, typs[0], typs[1:])
		}, elements...)
	case ast.Continue:
		recordBuiltin("Unit")
	case ast.NamedArgument:
		i.inferTypesOf(func(typs []types.Type) {
			i.iw.Record(element, typs[0].Symbol())
		}, n.Value())

	default:
		i.typeCache.error(element, "Inferring type for unknown node: %s", element)
	}
}

func (i *inferencer) inferTypesOf(pend func(types []types.Type), elements ...ast.Element) {
	waitingTypes := len(elements)
	result := make([]types.Type, waitingTypes)
	for index, element := range elements {
		i.Infer(element)
		i.iw.Pend(element, func(typSym types.TypeSymbol) {
			i.typeWorker.Pend(typSym, func() {
				result[index] = typSym.Type()
				waitingTypes--
				if waitingTypes == 0 {
					pend(result)
				}
			})
		})
	}
}

func (i *inferencer) resolveCall(call ast.Call, callType types.Type, argumentTypes []types.Type) {
	var namedArguments []ast.NamedArgument
	var namedArgumentTypes []types.Type
	var positionalArguments []ast.Element
	var positionalArgumentTypes []types.Type
	for i, argument := range call.Arguments() {
		switch n := argument.(type) {
		case ast.NamedArgument:
			namedArguments = append(namedArguments, n)
			namedArgumentTypes = append(namedArgumentTypes, argumentTypes[i])
		default:
			positionalArguments = append(positionalArguments, argument)
			positionalArgumentTypes = append(positionalArgumentTypes, argumentTypes[i])
		}
	}
	var candidate types.Signature
signatureLoop:
	for _, signature := range callType.Signatures() {
		var parameters = make(map[string]types.Parameter)
		for _, parameter := range signature.Parameters() {
			parameters[parameter.Name()] = parameter
		}

		usedParameters := make(map[types.Parameter]types.Parameter)

		// First select all named parameters
		for i, namedArgument := range namedArguments {
			namedArgumentType := namedArgumentTypes[i]
			namedParameter, ok := parameters[namedArgument.Name().Text()]
			if !ok {
				continue signatureLoop
			}
			if namedParameter.Type() != namedArgumentType {
				continue signatureLoop
			}
			usedParameters[namedParameter] = namedParameter
		}

		// Check the unused parameters
		argumentIndex := 0
		for _, parameter := range signature.Parameters() {
			_, ok := usedParameters[parameter]
			if ok {
				continue
			}
			if argumentIndex < len(positionalArguments) {
				if positionalArgumentTypes[argumentIndex] != parameter.Type() {
					continue signatureLoop
				}
			} else {
				// Default parameters should be added here
				continue signatureLoop
			}
		}
		candidate = signature
		break
	}
	if candidate == nil {
		i.typeCache.error(call, "Could not find a matching signature")
		i.iw.Record(call, i.typeCache.ErrorType())
	} else {
		i.iw.Record(call, candidate.Result().Symbol())
	}
}

type elementTypes struct {
	types map[ast.Element]types.TypeSymbol
}

func newElementTypes() *elementTypes {
	return &elementTypes{types: make(map[ast.Element]types.TypeSymbol)}
}

type inferenceWorker struct {
	elementTypes *elementTypes
	pending      map[ast.Element][]func(typSym types.TypeSymbol)
}

func newInferenceWorker(elementTypes *elementTypes) *inferenceWorker {
	return &inferenceWorker{
		elementTypes: elementTypes,
		pending:      make(map[ast.Element][]func(typSym types.TypeSymbol)),
	}
}

func (iw *inferenceWorker) Pend(element ast.Element, pend func(typSym types.TypeSymbol)) {
	typSym, ok := iw.elementTypes.types[element]
	if ok {
		pend(typSym)
	} else {
		pending, _ := iw.pending[element]
		iw.pending[element] = append(pending, pend)
	}
}

func (iw *inferenceWorker) Record(element ast.Element, typSym types.TypeSymbol) {
	_, ok := iw.elementTypes.types[element]
	assert.Assert(!ok, "AST typed multiple times: %s", element)
	iw.elementTypes.types[element] = typSym
	pending, ok := iw.pending[element]
	if ok {
		delete(iw.pending, element)
		for _, pend := range pending {
			pend(typSym)
		}
	}
}
