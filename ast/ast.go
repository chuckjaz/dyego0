package ast

import (
	"go/token"
)

// Element is the root of all AST nodes.
type Element interface {
	Locatable
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

// LiteralInt is a integer literalT
type LiteralInt interface {
	Element
	Value() int
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

// ObjectInitializer is an object intializer
type ObjectInitializer interface {
	Element
	Mutable() bool
	Members() []MemberInitializer
	IsObject() bool
}

// ArrayInitializer is an array initializer
type  ArrayInitializer interface {
	Element
        Mutable() bool
	Elements() []Element
	IsArray() bool
}

// MemberInitializer an abstraction for all object member initializers
type MemberInitializer interface {
	Element
}

// NamedMemberInitializer is a  named field initializer for an object
type NamedMemberInitializer interface {
	Element
	Name() Name
	Type() Element
	Value() Element
}

// SplatMemberInitializer is a type hoisting member intiailizer
type SplatMemberInitializer interface {
	MemberInitializer
	Type() Element
	IsSplat() bool
}

// Lambda is a lambda
type Lambda interface {
	Element
	Parameters() []Parameter
	Body() Element
}

// Parameter is a function or lambda parameter declaration
type Parameter interface {
	Element
	Name() Name
	Type() Element
	Default() Element
	IsParameter() bool
}

// VarDefinition is a field or local variable defintion or declaration
type VarDefinition interface {
	Element
	Name() Name
	Type() Element
	Mutable() bool
	IsField() bool
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
	UpdateContext()
	Name(value string) Name
	LiteralInt(value int) LiteralInt
	LiteralDouble(value float64) LiteralDouble
	LiteralFloat(value float32) LiteralFloat
	LiteralBoolean(value bool) LiteralBoolean
	LiteralString(value string) LiteralString
	LiteralNull() LiteralNull
	ObjectInitializer(mutable bool, members []MemberInitializer) ObjectInitializer
	ArrayInitializer(mutable bool, elements []Element) ArrayInitializer
	NamedMemberInitializer(name Name, typ Element, value Element) NamedMemberInitializer
	SplatMemberInitializer(typ Element) SplatMemberInitializer
	Lambda(parameters []Parameter, body Element) Lambda
	Parameter(name Name, typ Element, deflt Element) Parameter
	VarDefinition(name Name, typ Element, mutable bool) VarDefinition
	Error(message string) Error
}

// BuilderContext provides the source location for a Builder
type BuilderContext interface {
	Start() token.Pos
	End() token.Pos
}

type position struct {
	start, end token.Pos
}

type builderImpl struct {
	context BuilderContext
	locations []position
}

// NewBuilder makes an AST builder that can be used to make AST nodes
func NewBuilder(context BuilderContext) Builder {
	return &builderImpl{context: context}
}

func (b *builderImpl) PushContext() {
	b.locations = append(b.locations, position{start: b.context.Start(), end: b.context.End()})
}

func (b *builderImpl) PopContext() {
	l := len(b.locations)
	if l > 1 {
		last := b.locations[l-1]
		prev := b.locations[l-2]
		b.locations = append(b.locations[:l-2], position{start: prev.start, end: last.end})
	}
}

func (b *builderImpl) UpdateContext() {
	b.locations[len(b.locations)-1] = position{start: b.context.Start(), end: b.context.End()}
}

func (b *builderImpl) Loc() Location {
	pos := b.locations[len(b.locations)-1]
	return NewLocation(pos.start, pos.end)
}

type nameImpl struct {
	Location
	text string
}

func (n *nameImpl) Text() string {
	return n.text
}

func (b *builderImpl) Name(text string) Name {
	return &nameImpl{Location: b.Loc(), text: text}
}

type literalIntImpl struct {
	Location
	value int
}

func (l *literalIntImpl) Value() int {
	return l.value
}

func (b *builderImpl) LiteralInt(value int) LiteralInt {
	return &literalIntImpl{Location: b.Loc(), value: value}
}

type literalDoubleImpl struct {
	Location
	value float64
}

func (l *literalDoubleImpl) Value() float64 {
	return l.value
}

func (b *builderImpl) LiteralDouble(value float64) LiteralDouble {
	return &literalDoubleImpl{Location: b.Loc(), value: value}
}

type literalFloatImpl struct {
	Location
	value float32
}

func (l *literalFloatImpl) Value() float32 {
	return l.value
}

func (b *builderImpl) LiteralFloat(value float32) LiteralFloat {
	return &literalFloatImpl{Location: b.Loc(), value: value}
}

type literalBooleanImpl struct {
	Location
	value bool
}

func (l *literalBooleanImpl) Value() bool {
	return l.value
}

func (b *builderImpl) LiteralBoolean(value bool) LiteralBoolean {
	return &literalBooleanImpl{Location: b.Loc(), value: value}
}

type literalStringImpl struct {
	Location
	value string
}

func (l *literalStringImpl) Value() string {
	return l.value
}

func (b *builderImpl) LiteralString(value string) LiteralString {
	return &literalStringImpl{Location: b.Loc(), value: value}
}

type literalNullImpl struct {
	Location
}

func (l *literalNullImpl) IsNull() bool {
	return true
}

func (b *builderImpl) LiteralNull() LiteralNull {
	return &literalNullImpl{Location: b.Loc()}
}

type objectInitializerImpl struct {
	Location
	mutable bool
	members []MemberInitializer
}

func (o *objectInitializerImpl) Mutable() bool {
	return o.mutable
}

func (o *objectInitializerImpl) Members() []MemberInitializer {
	return o.members
}

func (o *objectInitializerImpl) IsObject() bool {
	return true
}

func (b *builderImpl) ObjectInitializer(mutable bool, members []MemberInitializer) ObjectInitializer {
	return &objectInitializerImpl{Location: b.Loc(), mutable: mutable, members: members}
}

type arrayInitializerImpl struct {
	Location
	mutable bool
	elements []Element
}

func (a *arrayInitializerImpl) Mutable() bool {
	return a.mutable
}

func (a *arrayInitializerImpl) Elements() []Element {
	return a.elements
}

func (a *arrayInitializerImpl) IsArray() bool {
	return true
}

func (b *builderImpl) ArrayInitializer(mutable bool, elements []Element) ArrayInitializer {
	return &arrayInitializerImpl{Location: b.Loc(), mutable: mutable, elements: elements}
}

type namedMemberInitializerImpl struct {
	Location
	name Name
	typ Element
	value Element
}

func (n *namedMemberInitializerImpl) Name() Name {
	return n.name
}

func (n *namedMemberInitializerImpl) Type() Element {
	return n.typ
}

func (n *namedMemberInitializerImpl) Value() Element {
	return n.value
}

func (b *builderImpl) NamedMemberInitializer(name Name, typ Element, value Element) NamedMemberInitializer {
	return &namedMemberInitializerImpl{Location: b.Loc(), name: name, typ: typ, value: value}
}

type splatMemberInitializerImpl struct {
	Location
	typ Element
}

func (s *splatMemberInitializerImpl) Type() Element {
	return s.typ
}

func (s *splatMemberInitializerImpl) IsSplat() bool {
	return true
}

func (b *builderImpl) SplatMemberInitializer(typ Element) SplatMemberInitializer {
	return &splatMemberInitializerImpl{Location: b.Loc(), typ: typ}
}

type lambdaImpl struct {
	Location
	parameters []Parameter
	body Element
}

func (l *lambdaImpl) Parameters() []Parameter {
	return l.parameters
}

func (l *lambdaImpl) Body() Element {
	return l.body
}

func (b *builderImpl) Lambda(parameters []Parameter, body Element) Lambda {
	return &lambdaImpl{Location: b.Loc(), parameters: parameters, body: body}
}

type parameterImpl struct {
	Location
	name Name
	typ Element
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

func (b *builderImpl) Parameter(name Name, typ Element, deflt Element) Parameter {
	return &parameterImpl{Location: b.Loc(), name: name, typ: typ, deflt: deflt}
}

type varDefinitionImpl struct {
	Location
	name Name
	typ Element
	mutable bool
}

func (v *varDefinitionImpl) Name() Name {
	return v.name
}

func (v *varDefinitionImpl) Type() Element {
	return v.typ
}

func (v *varDefinitionImpl) Mutable() bool {
	return v.mutable
}

func (v *varDefinitionImpl) IsField() bool {
	return true
}

func (b *builderImpl) VarDefinition(name Name, typ Element, mutable bool) VarDefinition {
	return &varDefinitionImpl{Location: b.Loc(), name: name, typ: typ, mutable: mutable}
}

type errorImpl struct {
	Location
	message string
}

func (e *errorImpl) Message() string {
	return e.message
}

func (b *builderImpl) Error(message string) Error {
	return &errorImpl{Location: b.Loc(), message: message}
}

