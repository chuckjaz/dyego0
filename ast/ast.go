package ast

import (
    "dyego0/location"
	"dyego0/tokens"
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

// LiteralRune is a rune literal
type LiteralRune interface {
	Element
	Value() rune
}

// LiteralByte is a byte literal
type LiteralByte interface {
	Element
	Value() byte
}

// LiteralInt is a integer literal
type LiteralInt interface {
	Element
	Value() int
}

// LiteralUInt is an unsigned integer literal
type LiteralUInt interface {
	Element
	Value() uint
}

// LiteralLong is a long literal
type LiteralLong interface {
	Element
	Value() int64
}

// LiteralULong is an unsigned long literal
type LiteralULong interface {
	Element
	Value() uint64
}

// LiteralDouble is a dobule literal
type LiteralDouble interface {
	Element
	Value() float64
}

// LiteralFloat is a float literal
type LiteralFloat interface {
	Element
	Value() float32
}

// LiteralString is a string literal
type LiteralString interface {
	Element
	Value() string
}

// LiteralBoolean is a boolean literal
type LiteralBoolean interface {
	Element
	Value() bool
}

// LiteralNull is a null litearl
type LiteralNull interface {
	Element
	IsNull() bool
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
}

// SpreadMemberInitializer is a type hoisting member intiailizer
type SpreadMemberInitializer interface {
	Element
	Target() Element
	IsSpreadMemberInitializer() bool
}

// Lambda is a lambda
type Lambda interface {
	Element
	TypeParameters() TypeParameters
	Parameters() []Parameter
	Body() Element
	IsLambda() bool
}

// IntrinsicLambda is that specifies instructions
type IntrinsicLambda interface {
	Element
	TypeParameters() TypeParameters
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

// TypeParameters is a type paramters clause
type TypeParameters interface {
	Element
	Parameters() []TypeParameter
	Wheres() []Where
}

// TypeParameter is a type parameter
type TypeParameter interface {
	Element
	Name() Name
	Constraint() Element
	IsTypeParameter() bool
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

// Where is a type parameter where clause
type Where interface {
	Element
	Left() Element
	Right() Element
	IsWhere() bool
}

// LetDefinition is a constant or type declaration
type LetDefinition interface {
	Element
	Name() Element
	Value() Element
	IsLetDefinition() bool
}

// VarDefinition is a field or local variable defintion or declaration
type VarDefinition interface {
	Element
	Name() Name
	Type() Element
	Value() Element
	Mutable() bool
	IsField() bool
}

// TypeLiteral is a type literal
type TypeLiteral interface {
	Element
	Members() []Element
	IsTypeLiteral() bool
}

// TypeLiteralMember is a type literal member
type TypeLiteralMember interface {
	Element
	Name() Name
	Type() Element
	IsTypeLiteralMember() bool
}

// TypeLiteralConstant is a type literal constant
type TypeLiteralConstant interface {
	Element
	Name() Name
	Value() Element
	IsTypeLiteralConstant() bool
}

// CallableTypeMember is a callable type literal member
type CallableTypeMember interface {
	Element
	Parameters() []Element
	Result() Element
}

// SpreadTypeMember is a spread of another type into a type literal
type SpreadTypeMember interface {
	Element
	Reference() Element
	IsSpreadTypeMember() bool
}

// SequenceType transforms Elements() type reference into a sequence
type SequenceType interface {
	Element
	Elements() Element
	IsSequenceType() bool
}

// OptionalType transforms Element() type reference into an obtional type
type OptionalType interface {
	Element
	Element() Element
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

// Error is a error node that should be reported
type Error interface {
	Element
	Message() string
}

// Builder is a helper to build AST nodes that uses context typecially provided by a scanner to
// initiazlier the location of AST nodes.
type Builder interface {
	PushContext()
	PopContext()
	Name(value string) Name
	LiteralRune(value rune) LiteralRune
	LiteralByte(value byte) LiteralByte
	LiteralInt(value int) LiteralInt
	LiteralUInt(value uint) LiteralUInt
	LiteralLong(value int64) LiteralLong
	LiteralULong(value uint64) LiteralULong
	LiteralDouble(value float64) LiteralDouble
	LiteralFloat(value float32) LiteralFloat
	LiteralBoolean(value bool) LiteralBoolean
	LiteralString(value string) LiteralString
	LiteralNull() LiteralNull
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
	SpreadMemberInitializer(target Element) SpreadMemberInitializer
	Lambda(typeParameters TypeParameters, parameters []Parameter, body Element) Lambda
	IntrinsicLambda(typeParameters TypeParameters, parameters []Parameter, body Element, result Element) IntrinsicLambda
	Loop(label Name, body Element) Loop
	Return(value Element) Return
	TypeParameters(parameters []TypeParameter, wheres []Where) TypeParameters
	TypeParameter(name Name, constraint Element) TypeParameter
	When(target Element, clauses []Element) When
	WhenValueClause(value Element, body Element) WhenValueClause
	WhenElseClause(body Element) WhenElseClause
	Where(left, right Element) Where
	Parameter(name Name, typ Element, deflt Element) Parameter
	VarDefinition(name Name, typ Element, value Element, mutable bool) VarDefinition
	LetDefinition(name Element, value Element) LetDefinition
	TypeLiteral(members []Element) TypeLiteral
	TypeLiteralConstant(name Name, value Element) TypeLiteralConstant
	TypeLiteralMember(name Name, typ Element) TypeLiteralMember
	CallableTypeMember(parameters []Element, result Element) CallableTypeMember
	SpreadTypeMember(reference Element) SpreadTypeMember
	SequenceType(elements Element) SequenceType
	OptionalType(element Element) OptionalType
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
	Error(message string) Error
	DirectError(start, end tokens.Pos, message string) Error
	Clone(context BuilderContext) Builder
	Loc() location.Location
}

// BuilderContext provides the source location for a Builder
type BuilderContext interface {
	Start() tokens.Pos
	End() tokens.Pos
}

type builderImpl struct {
	context   BuilderContext
	locations []tokens.Pos
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

type literalRuneImpl struct {
	location.Location
	value rune
}

func (l *literalRuneImpl) Value() rune {
	return l.value
}

func (l *literalRuneImpl) String() string {
	return fmt.Sprintf("LiteralRune(%s, '%s')", l.Location, string(l.value))
}

func (b *builderImpl) LiteralRune(value rune) LiteralRune {
	return &literalRuneImpl{Location: b.Loc(), value: value}
}

type literalByteImpl struct {
	location.Location
	value byte
}

func (l *literalByteImpl) Value() byte {
	return l.value
}

func (l *literalByteImpl) String() string {
	return fmt.Sprintf("LiteralByte(%s, %d)", l.Location, l.value)
}

func (b *builderImpl) LiteralByte(value byte) LiteralByte {
	return &literalByteImpl{Location: b.Loc(), value: value}
}

type literalIntImpl struct {
	location.Location
	value int
}

func (l *literalIntImpl) Value() int {
	return l.value
}

func (l *literalIntImpl) String() string {
	return fmt.Sprintf("LiteralInt(%s, %d)", l.Location, l.value)
}

func (b *builderImpl) LiteralInt(value int) LiteralInt {
	return &literalIntImpl{Location: b.Loc(), value: value}
}

type literalUIntImpl struct {
	location.Location
	value uint
}

func (l *literalUIntImpl) Value() uint {
	return l.value
}

func (l *literalUIntImpl) String() string {
	return fmt.Sprintf("LiteralUInt(%s, %d)", l.Location, l.value)
}

func (b *builderImpl) LiteralUInt(value uint) LiteralUInt {
	return &literalUIntImpl{Location: b.Loc(), value: value}
}

type literalLongImpl struct {
	location.Location
	value int64
}

func (l *literalLongImpl) Value() int64 {
	return l.value
}

func (l *literalLongImpl) String() string {
	return fmt.Sprintf("LiteralLong(%s, %d)", l.Location, l.value)
}

func (b *builderImpl) LiteralLong(value int64) LiteralLong {
	return &literalLongImpl{Location: b.Loc(), value: value}
}

type literalULongImpl struct {
	location.Location
	value uint64
}

func (l *literalULongImpl) Value() uint64 {
	return l.value
}

func (l *literalULongImpl) String() string {
	return fmt.Sprintf("LiteralULong(%s, %d)", l.Location, l.value)
}

func (b *builderImpl) LiteralULong(value uint64) LiteralULong {
	return &literalULongImpl{Location: b.Loc(), value: value}
}

type literalDoubleImpl struct {
	location.Location
	value float64
}

func (l *literalDoubleImpl) Value() float64 {
	return l.value
}

func (l *literalDoubleImpl) String() string {
	return fmt.Sprintf("LiteralDouble(%s, %v)", l.Location, l.value)
}

func (b *builderImpl) LiteralDouble(value float64) LiteralDouble {
	return &literalDoubleImpl{Location: b.Loc(), value: value}
}

type literalFloatImpl struct {
	location.Location
	value float32
}

func (l *literalFloatImpl) Value() float32 {
	return l.value
}

func (l *literalFloatImpl) String() string {
	return fmt.Sprintf("LiteralFloat(%s, %v)", l.Location, l.value)
}

func (b *builderImpl) LiteralFloat(value float32) LiteralFloat {
	return &literalFloatImpl{Location: b.Loc(), value: value}
}

type literalBooleanImpl struct {
	location.Location
	value bool
}

func (l *literalBooleanImpl) Value() bool {
	return l.value
}

func (l *literalBooleanImpl) String() string {
	return fmt.Sprintf("LiteralBoolean(%s, %v)", l.Location, l.value)
}

func (b *builderImpl) LiteralBoolean(value bool) LiteralBoolean {
	return &literalBooleanImpl{Location: b.Loc(), value: value}
}

type literalStringImpl struct {
	location.Location
	value string
}

func (l *literalStringImpl) Value() string {
	return l.value
}

func (l *literalStringImpl) String() string {
	return fmt.Sprintf("LiteralString(%s, \"%s\")", l.Location, l.value)
}

func (b *builderImpl) LiteralString(value string) LiteralString {
	return &literalStringImpl{Location: b.Loc(), value: value}
}

type literalNullImpl struct {
	location.Location
}

func (l *literalNullImpl) IsNull() bool {
	return true
}

func (l *literalNullImpl) String() string {
	return fmt.Sprintf("LiteralNull(%s)", l.Location)
}

func (b *builderImpl) LiteralNull() LiteralNull {
	return &literalNullImpl{Location: b.Loc()}
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

func (b *builderImpl) NamedMemberInitializer(name Name, typ Element, value Element) NamedMemberInitializer {
	return &namedMemberInitializerImpl{Location: b.Loc(), name: name, typ: typ, value: value}
}

type spreadMemberInitializerImpl struct {
	location.Location
	target Element
}

func (n *spreadMemberInitializerImpl) Target() Element {
	return n.target
}

func (n *spreadMemberInitializerImpl) IsSpreadMemberInitializer() bool {
	return true
}

func (n *spreadMemberInitializerImpl) String() string {
	return fmt.Sprintf("SpreadMemberInitializer(%s, target: %s)", n.Location, s(n.target))
}

func (b *builderImpl) SpreadMemberInitializer(target Element) SpreadMemberInitializer {
	return &spreadMemberInitializerImpl{Location: b.Loc(), target: target}
}

type lambdaImpl struct {
	location.Location
	typeParameters TypeParameters
	parameters     []Parameter
	body           Element
}

func (l *lambdaImpl) TypeParameters() TypeParameters {
	return l.typeParameters
}

func (l *lambdaImpl) Parameters() []Parameter {
	return l.parameters
}

func (l *lambdaImpl) Body() Element {
	return l.body
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
	return fmt.Sprintf("Lambda(%s, typeParameters: %s, parameters: %s, body: %s)", l.Location, s(l.typeParameters),
		parametersToString(l.parameters), s(l.body))
}

func (b *builderImpl) Lambda(typeParameters TypeParameters, parameters []Parameter, body Element) Lambda {
	return &lambdaImpl{Location: b.Loc(), typeParameters: typeParameters, parameters: parameters, body: body}
}

type intrinsicLambdaImpl struct {
	location.Location
	typeParameters TypeParameters
	parameters     []Parameter
	body           Element
	result         Element
}

func (l *intrinsicLambdaImpl) TypeParameters() TypeParameters {
	return l.typeParameters
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
	return fmt.Sprintf("IntrinsicLambda(%s, typeParameters: %s, parameters: %s, body: %s, result: %s)", l.Location, s(l.typeParameters),
		parametersToString(l.parameters), s(l.body), s(l.result))
}

func (b *builderImpl) IntrinsicLambda(typeParameters TypeParameters, parameters []Parameter, body Element, result Element) IntrinsicLambda {
	return &intrinsicLambdaImpl{Location: b.Loc(), typeParameters: typeParameters, parameters: parameters, body: body, result: result}
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

type typeParametersImpl struct {
	location.Location
	parameters []TypeParameter
	wheres     []Where
}

func (t *typeParametersImpl) Parameters() []TypeParameter {
	return t.parameters
}

func (t *typeParametersImpl) Wheres() []Where {
	return t.wheres
}

// LowerTypeParameters converts []TypeParameter to []Element
func LowerTypeParameters(parameters []TypeParameter) []Element {
	var result = make([]Element, len(parameters))
	for i, e := range parameters {
		result[i] = e
	}
	return result
}

// LowerWheres converst []Where to []Element
func LowerWheres(parameters []Where) []Element {
	var result = make([]Element, len(parameters))
	for i, e := range parameters {
		result[i] = e
	}
	return result
}

func (t *typeParametersImpl) String() string {
	return fmt.Sprintf("TypeParameters(%s, parameters: %s, wheres: %s)", t.Location,
		elementsToString(LowerTypeParameters(t.parameters)),
		elementsToString(LowerWheres(t.wheres)))
}

func (b *builderImpl) TypeParameters(parameters []TypeParameter, wheres []Where) TypeParameters {
	return &typeParametersImpl{Location: b.Loc(), parameters: parameters, wheres: wheres}
}

type typeParameterImpl struct {
	location.Location
	name       Name
	constraint Element
}

func (t *typeParameterImpl) Name() Name {
	return t.name
}

func (t *typeParameterImpl) Constraint() Element {
	return t.constraint
}

func (t *typeParameterImpl) IsTypeParameter() bool {
	return true
}

func (t *typeParameterImpl) String() string {
	return fmt.Sprintf("TypeParameter(%s, name: %s, constraint: %s)", t.Location, s(t.name), s(t.constraint))
}

func (b *builderImpl) TypeParameter(name Name, constraint Element) TypeParameter {
	return &typeParameterImpl{Location: b.Loc(), name: name, constraint: constraint}
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

type whereImpl struct {
	location.Location
	left  Element
	right Element
}

func (w *whereImpl) Left() Element {
	return w.left
}

func (w *whereImpl) Right() Element {
	return w.right
}

func (w *whereImpl) IsWhere() bool {
	return true
}

func (w *whereImpl) String() string {
	return fmt.Sprintf("Where(%s, left: %s, right: %s)", w.Location, s(w.left), s(w.right))
}

func (b *builderImpl) Where(left, right Element) Where {
	return &whereImpl{Location: b.Loc(), left: left, right: right}
}

type varDefinitionImpl struct {
	location.Location
	name    Name
	typ     Element
	value   Element
	mutable bool
}

func (v *varDefinitionImpl) Name() Name {
	return v.name
}

func (v *varDefinitionImpl) Type() Element {
	return v.typ
}

func (v *varDefinitionImpl) Value() Element {
	return v.value
}

func (v *varDefinitionImpl) Mutable() bool {
	return v.mutable
}

func (v *varDefinitionImpl) IsField() bool {
	return true
}

func (v *varDefinitionImpl) String() string {
	return fmt.Sprintf("VarDefinition(%s, name: %s, type: %s, value: %s, mutable: %v)", v.Location, s(v.name), s(v.typ), s(v.value),
		v.mutable)
}

func (b *builderImpl) VarDefinition(name Name, typ Element, value Element, mutable bool) VarDefinition {
	return &varDefinitionImpl{Location: b.Loc(), name: name, typ: typ, value: value, mutable: mutable}
}

type letDefinitionImpl struct {
	location.Location
	name  Element
	value Element
}

func (l *letDefinitionImpl) Name() Element {
	return l.name
}

func (l *letDefinitionImpl) Value() Element {
	return l.value
}

func (l *letDefinitionImpl) IsLetDefinition() bool {
	return true
}

func (l *letDefinitionImpl) String() string {
	return fmt.Sprintf("LetDefinition(%s, name: %s, value: %s)", l.Location, s(l.name), s(l.value))
}

func (b *builderImpl) LetDefinition(name Element, value Element) LetDefinition {
	return &letDefinitionImpl{Location: b.Loc(), name: name, value: value}
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

type typeLiteralConstantImpl struct {
	location.Location
	name  Name
	value Element
}

func (t *typeLiteralConstantImpl) Name() Name {
	return t.name
}

func (t *typeLiteralConstantImpl) Value() Element {
	return t.value
}

func (t *typeLiteralConstantImpl) IsTypeLiteralConstant() bool {
	return true
}

func (t *typeLiteralConstantImpl) String() string {
	return fmt.Sprintf("TypeLiteralConstant(%s, name: %s, value: %s)", t.Location, s(t.name), s(t.value))
}

func (b *builderImpl) TypeLiteralConstant(name Name, value Element) TypeLiteralConstant {
	return &typeLiteralConstantImpl{Location: b.Loc(), name: name, value: value}
}

type typeLiteralMemberImpl struct {
	location.Location
	name Name
	typ  Element
}

func (m *typeLiteralMemberImpl) Name() Name {
	return m.name
}

func (m *typeLiteralMemberImpl) Type() Element {
	return m.typ
}

func (m *typeLiteralMemberImpl) IsTypeLiteralMember() bool {
	return true
}

func (m *typeLiteralMemberImpl) String() string {
	return fmt.Sprintf("TypeLiteralMember(%s, name: %s, type: %s)", m.Location, s(m.name), s(m.typ))
}

func (b *builderImpl) TypeLiteralMember(name Name, typ Element) TypeLiteralMember {
	return &typeLiteralMemberImpl{Location: b.Loc(), name: name, typ: typ}
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

type spreadTypeMemberImpl struct {
	location.Location
	reference Element
}

func (m *spreadTypeMemberImpl) Reference() Element {
	return m.reference
}

func (m *spreadTypeMemberImpl) IsSpreadTypeMember() bool {
	return true
}

func (m *spreadTypeMemberImpl) String() string {
	return fmt.Sprintf("SpreadTypeMember(%s, reference: %s)", m.Location, s(m.reference))
}

func (b *builderImpl) SpreadTypeMember(reference Element) SpreadTypeMember {
	return &spreadTypeMemberImpl{Location: b.Loc(), reference: reference}
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

type optionalTypeImpl struct {
	location.Location
	element Element
}

func (o *optionalTypeImpl) Element() Element {
	return o.element
}

func (o *optionalTypeImpl) IsOptionalType() bool {
	return true
}

func (o *optionalTypeImpl) String() string {
	return fmt.Sprintf("OptionalType(%s, element: %s)", o.Location, s(o.element))
}

func (b *builderImpl) OptionalType(element Element) OptionalType {
	return &optionalTypeImpl{Location: b.Loc(), element: element}
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

type errorImpl struct {
	location.Location
	message string
}

func (e *errorImpl) Message() string {
	return e.message
}

func (e *errorImpl) String() string {
	return fmt.Sprintf("Error(%s, message: %s)", e.Location, e.message)
}

func (b *builderImpl) Error(message string) Error {
	return &errorImpl{Location: b.Loc(), message: message}
}

func (b *builderImpl) DirectError(start, end tokens.Pos, message string) Error {
	return &errorImpl{Location: location.NewLocation(start, end), message: message}
}

func (b *builderImpl) Clone(context BuilderContext) Builder {
	return &builderImpl{context: context, locations: b.locations}
}
