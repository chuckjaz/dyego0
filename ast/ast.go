package ast

import (
	"dyego0/errors"
	"dyego0/location"
	"fmt"
)

// Element is the root of all AST nodes.
type Element interface {
	location.Locatable
}

// Name is an identifier name
type Name interface {
	Element
	Text() string
}

// NamedElement is an abstraction for all named nodes
type NamedElement interface {
	Element
	Name() Name
}

// TypedElement is an abstraction for all typed elements
type TypedElement interface {
	Element
	Type() Element
}

// Literal is a rune literal
type Literal interface {
	Element
	Value() interface{}
}

// Selection is a member selector
type Selection interface {
	Element
	Target() Element
	Member() Name
}

// Sequence is a sequence of expressions
type Sequence interface {
	Element
	Left() Element
	Right() Element
	IsSequence() bool
}

// Spread is a spread expression prefix
type Spread interface {
	Element
	Target() Element
	IsSpread() bool
}

// Break is a break statement
type Break interface {
	Element
	Label() Name
	IsBreak() bool
}

// Call is a call expression
type Call interface {
	Element
	Target() Element
	Arguments() []Element
}

// Continue is a continue statement
type Continue interface {
	Element
	Label() Name
	IsContinue() bool
}

// NamedArgument is a named argument to a call expression
type NamedArgument interface {
	Element
	Name() Name
	Value() Element
	IsNamedArgument() bool
}

// ObjectInitializer is an object intializer
type ObjectInitializer interface {
	Element
	Mutable() bool
	Type() Element
	Members() []Element
	IsObject() bool
}

// ArrayInitializer is an array initializer
type ArrayInitializer interface {
	Element
	Mutable() bool
	Type() Element
	Elements() []Element
	IsArray() bool
}

// NamedMemberInitializer is a  named field initializer for an object
type NamedMemberInitializer interface {
	Element
	Name() Name
	Type() Element
	Value() Element
	IsNamedMemberInitializer() bool
}

// Lambda is a lambda
type Lambda interface {
	Element
	Parameters() []Parameter
	Body() Element
	Result() Element
	IsLambda() bool
}

// IntrinsicLambda is that specifies instructions
type IntrinsicLambda interface {
	Element
	Parameters() []Parameter
	Body() Element
	Result() Element
	IsIntrinsicLambda() bool
}

// Loop is a loop statement
type Loop interface {
	Element
	Label() Name
	Body() Element
	IsLoop() bool
}

// Parameter is a function or lambda parameter declaration
type Parameter interface {
	Element
	Name() Name
	Type() Element
	Default() Element
	IsParameter() bool
}

// Return is a return statement
type Return interface {
	Element
	Value() Element
	IsReturn() bool
}

// When is a when expression
type When interface {
	Element
	Target() Element
	Clauses() []Element
}

// WhenValueClause is a value clause of a when expression
type WhenValueClause interface {
	Element
	Value() Element
	Body() Element
	IsWhenValueClause() bool
}

// WhenElseClause is an else clause of a when expression
type WhenElseClause interface {
	Element
	Body() Element
	IsElse() bool
}

// Definition is a constant or type declaration
type Definition interface {
	Element
	Name() Element
	Type() Element
	Value() Element
	IsDefinition() bool
}

// Storage is a field or local variable defintion or declaration
type Storage interface {
	Element
	Name() Name
	Type() Element
	Value() Element
	Mutable() bool
}

// TypeLiteral a type literal
type TypeLiteral interface {
	Element
	Members() []Element
	IsTypeLiteral() bool
}

// CallableTypeMember is a callable type literal member
type CallableTypeMember interface {
	Element
	Parameters() []Element
	Result() Element
}

// SequenceType transforms Elements() type reference into a sequence
type SequenceType interface {
	Element
	Elements() Element
	IsSequenceType() bool
}

// ReferenceType transforms Referent() type reference into a reference
type ReferenceType interface {
	Element
	Referent() Element
}

// OptionalType transformls Target() type reference into an optional type
type OptionalType interface {
	Element
	Target() Element
	IsOptionalType() bool
}

// VocabularyLiteral is a vocabulary literal
type VocabularyLiteral interface {
	Element
	Members() []Element
	IsVocabularyLiteral() bool
}

// OperatorPlacement declares the placement of an operator
type OperatorPlacement int

const (
	// Infix is an infix operator placement
	Infix OperatorPlacement = iota

	// Prefix is a prefix operator placement
	Prefix

	// Postfix is a postfix operator placement
	Postfix

	// UnspecifiedPlacement indicates the placement of an operator was not specified
	UnspecifiedPlacement
)

func (p OperatorPlacement) String() string {
	switch p {
	case Infix:
		return "infix"
	case Prefix:
		return "prefix"
	case Postfix:
		return "postfix"
	default:
		return "invalid placement"
	}
}

// OperatorAssociativity is the associativity of an operator
type OperatorAssociativity int

const (
	// Left declares an operator to be left associative
	Left OperatorAssociativity = iota

	// Right declares an operator to be right associative
	Right

	// UnspecifiedAssociativity is the value of associativity when it was not specified
	UnspecifiedAssociativity
)

func (a OperatorAssociativity) String() string {
	switch a {
	case Left:
		return "left"
	case Right:
		return "right"
	default:
		return "invalid associativity"
	}
}

// OperatorPrecedenceRelation defines a partial ordering with other operators
type OperatorPrecedenceRelation int

const (
	// Before indicates that the precedence is before the referenced operator
	Before OperatorPrecedenceRelation = iota

	// After indicates that the precedence is after the referenced operator
	After
)

func (r OperatorPrecedenceRelation) String() string {
	switch r {
	case Before:
		return "before"
	case After:
		return "after"
	default:
		return "invalid relation"
	}
}

// VocabularyOperatorDeclaration is the declaration of an operator
type VocabularyOperatorDeclaration interface {
	Element
	Names() []Name
	Placement() OperatorPlacement
	Precedence() VocabularyOperatorPrecedence
	Associativity() OperatorAssociativity
}

// VocabularyOperatorPrecedence is the definition of an operator precedence
type VocabularyOperatorPrecedence interface {
	Element
	Name() Name
	Placement() OperatorPlacement
	Relation() OperatorPrecedenceRelation
}

// VocabularyEmbedding is an embedding of a vocabulary
type VocabularyEmbedding interface {
	Element
	Name() []Name
	IsVocabularyEmbedding() bool
}

// Builder is a helper to build AST nodes that uses context typecially provided by a scanner to
// initiazlier the location of AST nodes.
type Builder interface {
	PushContext()
	PopContext()
	Name(value string) Name
	Literal(value interface{}) Literal
	Break(label Name) Break
	Continue(label Name) Continue
	Sequence(left, right Element) Sequence
	Selection(target Element, member Name) Selection
	Spread(target Element) Spread
	Call(target Element, arguments []Element) Call
	NamedArgument(name Name, value Element) NamedArgument
	ObjectInitializer(mutable bool, typ Element, members []Element) ObjectInitializer
	ArrayInitializer(mutable bool, typ Element, elements []Element) ArrayInitializer
	NamedMemberInitializer(name Name, typ Element, value Element) NamedMemberInitializer
	Lambda(parameters []Parameter, body Element, result Element) Lambda
	IntrinsicLambda(parameters []Parameter, body Element, result Element) IntrinsicLambda
	Loop(label Name, body Element) Loop
	Return(value Element) Return
	When(target Element, clauses []Element) When
	WhenValueClause(value Element, body Element) WhenValueClause
	WhenElseClause(body Element) WhenElseClause
	Parameter(name Name, typ Element, deflt Element) Parameter
	Definition(name Element, typ Element, value Element) Definition
	Storage(name Name, typ Element, value Element, mutable bool) Storage
	TypeLiteral(members []Element) TypeLiteral
	CallableTypeMember(parameters []Element, result Element) CallableTypeMember
	SequenceType(elements Element) SequenceType
	OptionalType(target Element) OptionalType
	ReferenceType(referent Element) ReferenceType
	VocabularyLiteral(members []Element) VocabularyLiteral
	VocabularyOperatorDeclaration(
		names []Name,
		placement OperatorPlacement,
		precedence VocabularyOperatorPrecedence,
		associativity OperatorAssociativity,
	) VocabularyOperatorDeclaration
	VocabularyOperatorPrecedence(
		name Name,
		placement OperatorPlacement,
		relation OperatorPrecedenceRelation,
	) VocabularyOperatorPrecedence
	VocabularyEmbedding(name []Name) VocabularyEmbedding
	Error(message string, args ...interface{}) errors.Error
	Clone(context BuilderContext) Builder
	Loc() location.Location
}

// BuilderContext provides the source location for a Builder
type BuilderContext interface {
	Start() location.Pos
	End() location.Pos
}

type builderImpl struct {
	context   BuilderContext
	locations []location.Pos
}

// NewBuilder makes an AST builder that can be used to make AST nodes
func NewBuilder(context BuilderContext) Builder {
	return &builderImpl{context: context}
}

func (b *builderImpl) PushContext() {
	b.locations = append(b.locations, b.context.Start())
}

func (b *builderImpl) PopContext() {
	b.locations = b.locations[0 : len(b.locations)-1]
}

func (b *builderImpl) Loc() location.Location {
	start := b.locations[len(b.locations)-1]
	return location.NewLocation(start, b.context.End())
}

type nameImpl struct {
	location.Location
	text string
}

func (n *nameImpl) Text() string {
	return n.text
}

func (n *nameImpl) String() string {
	return fmt.Sprintf("Name(%s, %s)", n.Location, n.text)
}

func (b *builderImpl) Name(text string) Name {
	return &nameImpl{Location: b.Loc(), text: text}
}

type literalImpl struct {
	location.Location
	value interface{}
}

func (l *literalImpl) Value() interface{} {
	return l.value
}

func sOf(value interface{}) string {
	switch v := value.(type) {
	case rune:
		return fmt.Sprintf("'%c'", v)
	case byte:
		return fmt.Sprintf("%db", v)
	case int:
		return fmt.Sprintf("%d", v)
	case uint:
		return fmt.Sprintf("%du", v)
	case int64:
		return fmt.Sprintf("%dl", v)
	case float32:
		return fmt.Sprintf("%ff", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case string:
		return fmt.Sprintf("\"%s\"", v)
	}
	return fmt.Sprintf("%s", value)
}

func (l *literalImpl) String() string {
	return fmt.Sprintf("Literal(%s, %s)", l.Location, sOf(l.value))
}

func (b *builderImpl) Literal(value interface{}) Literal {
	return &literalImpl{Location: b.Loc(), value: value}
}

type breakImpl struct {
	location.Location
	label Name
}

func (b *breakImpl) Label() Name {
	return b.label
}

func (b *breakImpl) IsBreak() bool {
	return true
}

func (b *breakImpl) String() string {
	return fmt.Sprintf("Break(%s)", b.Location)
}

func (b *builderImpl) Break(label Name) Break {
	return &breakImpl{Location: b.Loc(), label: label}
}

type continueImpl struct {
	location.Location
	label Name
}

func (b *continueImpl) Label() Name {
	return b.label
}

func (b *continueImpl) IsContinue() bool {
	return true
}

func (b *continueImpl) String() string {
	return fmt.Sprintf("Continue(%s)", b.Location)
}

func (b *builderImpl) Continue(label Name) Continue {
	return &continueImpl{Location: b.Loc(), label: label}
}

type selectionImpl struct {
	location.Location
	target Element
	member Name
}

func (l *selectionImpl) Target() Element {
	return l.target
}

func (l *selectionImpl) Member() Name {
	return l.member
}

func s(e Element) string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("%s", e)
}

func (l *selectionImpl) String() string {
	return fmt.Sprintf("Selection(%s, target: %s, member: %s)", l.Location, s(l.target), s(l.member))
}

func (b *builderImpl) Selection(target Element, member Name) Selection {
	return &selectionImpl{Location: b.Loc(), target: target, member: member}
}

type sequenceImpl struct {
	location.Location
	left  Element
	right Element
}

func (n *sequenceImpl) Left() Element {
	return n.left
}

func (n *sequenceImpl) Right() Element {
	return n.right
}

func (n *sequenceImpl) IsSequence() bool {
	return true
}

func (n *sequenceImpl) String() string {
	return fmt.Sprintf("Sequence(%s, left: %s, right: %s)", n.Location, s(n.left), s(n.right))
}

func (b *builderImpl) Sequence(left, right Element) Sequence {
	return &sequenceImpl{Location: b.Loc(), left: left, right: right}
}

type optionalTypeImpl struct {
	location.Location
	target Element
}

func (n *optionalTypeImpl) Target() Element {
	return n.target
}

func (n *optionalTypeImpl) String() string {
	return fmt.Sprintf("OptionalType(%s, target: %s", n.Location, s(n.target))
}

func (n *optionalTypeImpl) IsOptionalType() bool {
	return true
}

func (b *builderImpl) OptionalType(target Element) OptionalType {
	return &optionalTypeImpl{Location: b.Loc(), target: target}
}

type referenceImpl struct {
	location.Location
	referent Element
}

func (n *referenceImpl) Referent() Element {
	return n.referent
}

func (n *referenceImpl) String() string {
	return fmt.Sprintf("ReferenceType(%s, referent: %s)", n.Location, s(n.referent))
}

func (b *builderImpl) ReferenceType(referent Element) ReferenceType {
	return &referenceImpl{Location: b.Loc(), referent: referent}
}

type spreadImpl struct {
	location.Location
	target Element
}

func (n *spreadImpl) Target() Element {
	return n.target
}

func (n *spreadImpl) IsSpread() bool {
	return true
}

func (n *spreadImpl) String() string {
	return fmt.Sprintf("Spread(%s, target: %s)", n.Location, s(n.target))
}

func (b *builderImpl) Spread(target Element) Spread {
	return &spreadImpl{Location: b.Loc(), target: target}
}

type callImpl struct {
	location.Location
	target    Element
	arguments []Element
}

func (c *callImpl) Target() Element {
	return c.target
}

func (c *callImpl) Arguments() []Element {
	return c.arguments
}

func elementsToString(elements []Element) string {
	result := "["
	first := true
	for _, element := range elements {
		if !first {
			result += ", "
		}
		first = false
		result += fmt.Sprintf("%s", element)
	}
	result += "]"
	return result
}

func (c *callImpl) String() string {
	return fmt.Sprintf("Call(%s, target: %s, arguments: %s)", c.Location, s(c.target), elementsToString(c.arguments))
}

func (b *builderImpl) Call(target Element, arguments []Element) Call {
	return &callImpl{Location: b.Loc(), target: target, arguments: arguments}
}

type namedArgumentImpl struct {
	location.Location
	name  Name
	value Element
}

func (n *namedArgumentImpl) Name() Name {
	return n.name
}

func (n *namedArgumentImpl) Value() Element {
	return n.value
}

func (n *namedArgumentImpl) IsNamedArgument() bool {
	return true
}

func (n *namedArgumentImpl) String() string {
	return fmt.Sprintf("NamedArgument(%s, name: %s, value: %s)", n.Location, s(n.name), s(n.value))
}

func (b *builderImpl) NamedArgument(name Name, value Element) NamedArgument {
	return &namedArgumentImpl{Location: b.Loc(), name: name, value: value}
}

type objectInitializerImpl struct {
	location.Location
	mutable bool
	typ     Element
	members []Element
}

func (o *objectInitializerImpl) Mutable() bool {
	return o.mutable
}

func (o *objectInitializerImpl) Type() Element {
	return o.typ
}

func (o *objectInitializerImpl) Members() []Element {
	return o.members
}

func (o *objectInitializerImpl) IsObject() bool {
	return true
}

func (o *objectInitializerImpl) String() string {
	return fmt.Sprintf("ObjectInitializer(%s, mutable: %v, type: %s, members: %s)", o.Location, o.mutable,
		s(o.typ), elementsToString(o.members))
}

func (b *builderImpl) ObjectInitializer(mutable bool, typ Element, members []Element) ObjectInitializer {
	return &objectInitializerImpl{Location: b.Loc(), mutable: mutable, typ: typ, members: members}
}

type arrayInitializerImpl struct {
	location.Location
	mutable  bool
	typ      Element
	elements []Element
}

func (a *arrayInitializerImpl) Mutable() bool {
	return a.mutable
}

func (a *arrayInitializerImpl) Type() Element {
	return a.typ
}

func (a *arrayInitializerImpl) Elements() []Element {
	return a.elements
}

func (a *arrayInitializerImpl) IsArray() bool {
	return true
}

func (a *arrayInitializerImpl) String() string {
	return fmt.Sprintf("ArrayInitializer(%s, mutable: %v, type: %s, elements: %s)", a.Location, a.mutable,
		s(a.typ), elementsToString(a.elements))
}

func (b *builderImpl) ArrayInitializer(mutable bool, typ Element, elements []Element) ArrayInitializer {
	return &arrayInitializerImpl{Location: b.Loc(), typ: typ, mutable: mutable, elements: elements}
}

type namedMemberInitializerImpl struct {
	location.Location
	name  Name
	typ   Element
	value Element
}

func (n *namedMemberInitializerImpl) Name() Name {
	return n.name
}

func (n *namedMemberInitializerImpl) Value() Element {
	return n.value
}

func (n *namedMemberInitializerImpl) Type() Element {
	return n.typ
}

func (n *namedMemberInitializerImpl) String() string {
	return fmt.Sprintf("NamedMemberInitializer(%s, name: %s, type: %s, value: %s)", n.Location, s(n.name),
		s(n.typ), s(n.value))
}

func (n *namedMemberInitializerImpl) IsNamedMemberInitializer() bool {
	return true
}

func (b *builderImpl) NamedMemberInitializer(name Name, typ Element, value Element) NamedMemberInitializer {
	return &namedMemberInitializerImpl{Location: b.Loc(), name: name, typ: typ, value: value}
}

type lambdaImpl struct {
	location.Location
	parameters []Parameter
	body       Element
	result     Element
}

func (l *lambdaImpl) Parameters() []Parameter {
	return l.parameters
}

func (l *lambdaImpl) Body() Element {
	return l.body
}

func (l *lambdaImpl) Result() Element {
	return l.result
}

func (l *lambdaImpl) IsLambda() bool {
	return true
}

// LowerParameters converts a []Parameter to []Element
func LowerParameters(parameters []Parameter) []Element {
	var result = make([]Element, len(parameters))
	for i, e := range parameters {
		result[i] = e
	}
	return result
}

func parametersToString(parameters []Parameter) string {
	return elementsToString(LowerParameters(parameters))
}

func (l *lambdaImpl) String() string {
	return fmt.Sprintf("Lambda(%s, parameters: %s, body: %s, result: %s)", l.Location,
		parametersToString(l.parameters), s(l.body), s(l.result))
}

func (b *builderImpl) Lambda(parameters []Parameter, body Element, ret Element) Lambda {
	return &lambdaImpl{Location: b.Loc(), parameters: parameters, body: body}
}

type intrinsicLambdaImpl struct {
	location.Location
	parameters []Parameter
	body       Element
	result     Element
}

func (l *intrinsicLambdaImpl) Parameters() []Parameter {
	return l.parameters
}

func (l *intrinsicLambdaImpl) Body() Element {
	return l.body
}

func (l *intrinsicLambdaImpl) Result() Element {
	return l.result
}

func (l *intrinsicLambdaImpl) IsIntrinsicLambda() bool {
	return true
}

func (l *intrinsicLambdaImpl) String() string {
	return fmt.Sprintf("IntrinsicLambda(%s, parameters: %s, body: %s, result: %s)", l.Location,
		parametersToString(l.parameters), s(l.body), s(l.result))
}

func (b *builderImpl) IntrinsicLambda(parameters []Parameter, body Element, result Element) IntrinsicLambda {
	return &intrinsicLambdaImpl{Location: b.Loc(), parameters: parameters, body: body, result: result}
}

type loopImpl struct {
	location.Location
	label Name
	body  Element
}

func (l *loopImpl) Label() Name {
	return l.label
}

func (l *loopImpl) Body() Element {
	return l.body
}

func (l *loopImpl) IsLoop() bool {
	return true
}

func (l *loopImpl) String() string {
	return fmt.Sprintf("Loop(%s, label: %s, body: %s)", l.Location, s(l.label), s(l.body))
}

func (b *builderImpl) Loop(label Name, body Element) Loop {
	return &loopImpl{Location: b.Loc(), label: label, body: body}
}

type parameterImpl struct {
	location.Location
	name  Name
	typ   Element
	deflt Element
}

func (p *parameterImpl) Name() Name {
	return p.name
}

func (p *parameterImpl) Type() Element {
	return p.typ
}

func (p *parameterImpl) Default() Element {
	return p.deflt
}

func (p *parameterImpl) IsParameter() bool {
	return true
}

func (p *parameterImpl) String() string {
	return fmt.Sprintf("Parameter(%s, name: %s, type: %s, default: %s)", p.Location, s(p.name), s(p.typ), s(p.deflt))
}

func (b *builderImpl) Parameter(name Name, typ Element, deflt Element) Parameter {
	return &parameterImpl{Location: b.Loc(), name: name, typ: typ, deflt: deflt}
}

type returnImpl struct {
	location.Location
	value Element
}

func (r *returnImpl) Value() Element {
	return r.value
}

func (r *returnImpl) IsReturn() bool {
	return true
}

func (r *returnImpl) String() string {
	return fmt.Sprintf("Return(%s, value: %s)", r.Location, s(r.value))
}

func (b *builderImpl) Return(value Element) Return {
	return &returnImpl{Location: b.Loc(), value: value}
}

type whenImpl struct {
	location.Location
	target  Element
	clauses []Element
}

func (w *whenImpl) Target() Element {
	return w.target
}

func (w *whenImpl) Clauses() []Element {
	return w.clauses
}

func (w *whenImpl) String() string {
	return fmt.Sprintf("When(%s, target: %s, clauses: %s)", w.Location, s(w.target), elementsToString(w.clauses))
}

func (b *builderImpl) When(target Element, clauses []Element) When {
	return &whenImpl{Location: b.Loc(), target: target, clauses: clauses}
}

type whenElseClauseImpl struct {
	location.Location
	body Element
}

func (w *whenElseClauseImpl) Body() Element {
	return w.body
}

func (w *whenElseClauseImpl) IsElse() bool {
	return true
}

func (w *whenElseClauseImpl) String() string {
	return fmt.Sprintf("WhenElseClause(%s, body: %s)", w.Location, s(w.body))
}

func (b *builderImpl) WhenElseClause(body Element) WhenElseClause {
	return &whenElseClauseImpl{Location: b.Loc(), body: body}
}

type whenValueClauseImpl struct {
	location.Location
	value Element
	body  Element
}

func (w *whenValueClauseImpl) Value() Element {
	return w.value
}

func (w *whenValueClauseImpl) Body() Element {
	return w.body
}

func (w *whenValueClauseImpl) IsWhenValueClause() bool {
	return true
}

func (w *whenValueClauseImpl) String() string {
	return fmt.Sprintf("WhenValueClause(%s, value: %s, body: %s)", w.Location, s(w.value), s(w.body))
}

func (b *builderImpl) WhenValueClause(value Element, body Element) WhenValueClause {
	return &whenValueClauseImpl{Location: b.Loc(), value: value, body: body}
}

type storageImpl struct {
	location.Location
	name    Name
	typ     Element
	value   Element
	mutable bool
}

func (v *storageImpl) Name() Name {
	return v.name
}

func (v *storageImpl) Type() Element {
	return v.typ
}

func (v *storageImpl) Value() Element {
	return v.value
}

func (v *storageImpl) Mutable() bool {
	return v.mutable
}

func (v *storageImpl) String() string {
	return fmt.Sprintf("Storage(%s, name: %s, type: %s, value: %s, mutable: %v)", v.Location, s(v.name), s(v.typ), s(v.value),
		v.mutable)
}

func (b *builderImpl) Storage(name Name, typ Element, value Element, mutable bool) Storage {
	return &storageImpl{Location: b.Loc(), name: name, typ: typ, value: value, mutable: mutable}
}

type definitionImpl struct {
	location.Location
	name  Element
	typ   Element
	value Element
}

func (l *definitionImpl) Name() Element {
	return l.name
}

func (l *definitionImpl) Type() Element {
	return l.typ
}

func (l *definitionImpl) Value() Element {
	return l.value
}

func (l *definitionImpl) IsDefinition() bool {
	return true
}

func (l *definitionImpl) String() string {
	return fmt.Sprintf("Definition(%s, name: %s, type: %s, value: %s)", l.Location, s(l.name), s(l.typ), s(l.value))
}

func (b *builderImpl) Definition(name, typ, value Element) Definition {
	return &definitionImpl{Location: b.Loc(), name: name, typ: typ, value: value}
}

type typeLiteralImpl struct {
	location.Location
	members []Element
}

func (t *typeLiteralImpl) Members() []Element {
	return t.members
}

func (t *typeLiteralImpl) IsTypeLiteral() bool {
	return true
}

func (t *typeLiteralImpl) String() string {
	return fmt.Sprintf("TypeLiteral(%s, members: %s)", t.Location, elementsToString(t.members))
}

func (b *builderImpl) TypeLiteral(members []Element) TypeLiteral {
	return &typeLiteralImpl{Location: b.Loc(), members: members}
}

type callableTypeMemberImpl struct {
	location.Location
	parameters []Element
	result     Element
}

func (c *callableTypeMemberImpl) Parameters() []Element {
	return c.parameters
}

func (c *callableTypeMemberImpl) Result() Element {
	return c.result
}

func (c *callableTypeMemberImpl) String() string {
	return fmt.Sprintf("CallableTypeMember(%s, parameters: %s, result: %s)", c.Location, elementsToString(c.parameters), s(c.result))
}

func (b *builderImpl) CallableTypeMember(parameters []Element, result Element) CallableTypeMember {
	return &callableTypeMemberImpl{Location: b.Loc(), parameters: parameters, result: result}
}

type sequenceTypeImpl struct {
	location.Location
	elements Element
}

func (n *sequenceTypeImpl) Elements() Element {
	return n.elements
}

func (n *sequenceTypeImpl) IsSequenceType() bool {
	return true
}

func (n *sequenceTypeImpl) String() string {
	return fmt.Sprintf("SequenceType(%s, elements: %s)", n.Location, s(n.elements))
}

func (b *builderImpl) SequenceType(elements Element) SequenceType {
	return &sequenceTypeImpl{Location: b.Loc(), elements: elements}
}

type vocabularyLiteralImpl struct {
	location.Location
	members []Element
}

func (v *vocabularyLiteralImpl) Members() []Element {
	return v.members
}

func (v *vocabularyLiteralImpl) IsVocabularyLiteral() bool {
	return true
}

func (v *vocabularyLiteralImpl) String() string {
	return fmt.Sprintf("VocabularyLiteral(%s, members: %s)", v.Location, elementsToString(v.members))
}

func (b *builderImpl) VocabularyLiteral(members []Element) VocabularyLiteral {
	return &vocabularyLiteralImpl{Location: b.Loc(), members: members}
}

type vocabularyOperatorDeclarationImpl struct {
	location.Location
	names         []Name
	placement     OperatorPlacement
	precedence    VocabularyOperatorPrecedence
	associativity OperatorAssociativity
}

func (v *vocabularyOperatorDeclarationImpl) Names() []Name {
	return v.names
}

func (v *vocabularyOperatorDeclarationImpl) Placement() OperatorPlacement {
	return v.placement
}

func (v *vocabularyOperatorDeclarationImpl) Precedence() VocabularyOperatorPrecedence {
	return v.precedence
}

func (v *vocabularyOperatorDeclarationImpl) Associativity() OperatorAssociativity {
	return v.associativity
}

// LowerNames convers []Name to []Element
func LowerNames(a []Name) []Element {
	var result = make([]Element, len(a))
	for i, e := range a {
		result[i] = e
	}
	return result
}

func (v *vocabularyOperatorDeclarationImpl) String() string {
	return fmt.Sprintf("VocabularyOperatorDeclaration(%s, names: %s, placement: %s, precedence: %s, associativity: %s)",
		v.Location, elementsToString(LowerNames(v.names)), v.placement, s(v.precedence), v.associativity)
}

func (b *builderImpl) VocabularyOperatorDeclaration(
	names []Name,
	placement OperatorPlacement,
	precedence VocabularyOperatorPrecedence,
	associativity OperatorAssociativity,
) VocabularyOperatorDeclaration {
	return &vocabularyOperatorDeclarationImpl{
		Location:      b.Loc(),
		names:         names,
		placement:     placement,
		precedence:    precedence,
		associativity: associativity,
	}
}

type vocabularyOperatorPrecedenceImpl struct {
	location.Location
	name      Name
	placement OperatorPlacement
	relation  OperatorPrecedenceRelation
}

func (v *vocabularyOperatorPrecedenceImpl) Name() Name {
	return v.name
}

func (v *vocabularyOperatorPrecedenceImpl) Placement() OperatorPlacement {
	return v.placement
}

func (v *vocabularyOperatorPrecedenceImpl) Relation() OperatorPrecedenceRelation {
	return v.relation
}

func (v *vocabularyOperatorPrecedenceImpl) String() string {
	return fmt.Sprintf("VocabularyOperatorPrecedence(%s, name: %s, placement: %s, relation: %s)",
		v.Location, s(v.name), v.placement, v.relation)
}

func (b *builderImpl) VocabularyOperatorPrecedence(
	name Name,
	placement OperatorPlacement,
	relation OperatorPrecedenceRelation,
) VocabularyOperatorPrecedence {
	return &vocabularyOperatorPrecedenceImpl{
		Location:  b.Loc(),
		name:      name,
		placement: placement,
		relation:  relation,
	}
}

type vocabularyEmbeddingImpl struct {
	location.Location
	name []Name
}

func (v *vocabularyEmbeddingImpl) Name() []Name {
	return v.name
}

func (v *vocabularyEmbeddingImpl) IsVocabularyEmbedding() bool {
	return true
}

func (v *vocabularyEmbeddingImpl) String() string {
	return fmt.Sprintf("VocabularyEmbedding(%s, name: %s)", v.Location, elementsToString(LowerNames(v.name)))
}

func (b *builderImpl) VocabularyEmbedding(name []Name) VocabularyEmbedding {
	return &vocabularyEmbeddingImpl{Location: b.Loc(), name: name}
}

func (b *builderImpl) Error(message string, args ...interface{}) errors.Error {
	return errors.New(b.Loc(), message, args...)
}

func (b *builderImpl) Clone(context BuilderContext) Builder {
	return &builderImpl{context: context, locations: b.locations}
}
