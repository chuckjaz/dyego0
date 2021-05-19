package ast

import (
	"dyego0/assert"
	"dyego0/errors"
)

// Visitor is an AST visitor
type Visitor interface {
	Visit(element Element) bool
}

// Walk will call Visit with the element followed by walking the children
func Walk(element Element, visitor Visitor) bool {
	if element != nil {
		if visitor.Visit(element) {
			return WalkChildren(element, visitor)
		}
		return false
	}
	return true
}

// WalkChildren will walk the children of element
func WalkChildren(element Element, visitor Visitor) bool {
	for {
		switch e := element.(type) {
		case Name:
			return true
		case Literal:
			return true
		case Selection:
			return Walk(e.Target(), visitor) && Walk(e.Member(), visitor)
		case Sequence:
			return Walk(e.Left(), visitor) && Walk(e.Right(), visitor)
		case Spread:
			return Walk(e.Target(), visitor)
		case Break:
			return Walk(e.Label(), visitor)
		case Call:
			return Walk(e.Target(), visitor) && walkElements(e.Arguments(), visitor)
		case Continue:
			return Walk(e.Label(), visitor)
		case NamedArgument:
			return Walk(e.Name(), visitor) && Walk(e.Value(), visitor)
		case ObjectInitializer:
			return Walk(e.Type(), visitor) && walkElements(e.Members(), visitor)
		case ArrayInitializer:
			return Walk(e.Type(), visitor) && walkElements(e.Elements(), visitor)
		case NamedMemberInitializer:
			return Walk(e.Name(), visitor) && Walk(e.Type(), visitor) && Walk(e.Value(), visitor)
		case Lambda:
			return walkParameters(e.Parameters(), visitor) && Walk(e.Body(), visitor) &&
				Walk(e.Result(), visitor)
		case IntrinsicLambda:
			return walkParameters(e.Parameters(), visitor) && Walk(e.Body(), visitor) &&
				Walk(e.Result(), visitor)
		case Loop:
			return Walk(e.Label(), visitor) && Walk(e.Body(), visitor)
		case Parameter:
			return Walk(e.Name(), visitor) && Walk(e.Type(), visitor) && Walk(e.Default(), visitor)
		case Return:
			return Walk(e.Value(), visitor)
		case When:
			return Walk(e.Target(), visitor) && walkElements(e.Clauses(), visitor)
		case WhenValueClause:
			return Walk(e.Value(), visitor) && Walk(e.Body(), visitor)
		case WhenElseClause:
			return Walk(e.Body(), visitor)
		case Definition:
			return Walk(e.Name(), visitor) && Walk(e.Type(), visitor) && Walk(e.Value(), visitor)
		case Storage:
			return Walk(e.Name(), visitor) && Walk(e.Type(), visitor) && Walk(e.Value(), visitor)
		case TypeLiteral:
			return walkElements(e.Members(), visitor)
		case TypeLiteralMember:
			return Walk(e.Name(), visitor) && Walk(e.Type(), visitor)
		case TypeLiteralConstant:
			return Walk(e.Name(), visitor) && Walk(e.Value(), visitor)
		case CallableTypeMember:
			return walkElements(e.Parameters(), visitor) && Walk(e.Result(), visitor)
		case SequenceType:
			return Walk(e.Elements(), visitor)
		case ReferenceType:
			return Walk(e.Referent(), visitor)
		case VocabularyLiteral:
			return walkElements(e.Members(), visitor)
		case VocabularyOperatorDeclaration:
			return walkNames(e.Names(), visitor) && Walk(e.Precedence(), visitor)
		case VocabularyOperatorPrecedence:
			return Walk(e.Name(), visitor)
		case VocabularyEmbedding:
			return walkNames(e.Name(), visitor)
		case errors.Error:
			return true
		default:
			assert.Fail("Unknown element %#v", element)
			return false
		}
	}
}

func walkElements(elements []Element, visitor Visitor) bool {
	for _, element := range elements {
		if !Walk(element, visitor) {
			return false
		}
	}
	return true
}

func walkParameters(parameters []Parameter, visitor Visitor) bool {
	for _, parameter := range parameters {
		if !Walk(parameter, visitor) {
			return false
		}
	}
	return true
}

func walkNames(names []Name, visitor Visitor) bool {
	for _, name := range names {
		if !Walk(name, visitor) {
			return false
		}
	}
	return true
}
