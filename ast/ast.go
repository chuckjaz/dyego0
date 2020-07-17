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

// Call is a call expression
type Call interface {
    Element
    Target() Element
    Arguments() []Element
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
    TypeParameters() TypeParameters
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

// Where is a type parameter where clause
type Where interface {
    Element
    Left() Element
    Right() Element
    IsWhere() bool
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
    Selection(target Element, member Name) Selection
    Call(target Element, arguments []Element) Call
    NamedArgument(name Name, value Element) NamedArgument
    ObjectInitializer(mutable bool, members []MemberInitializer) ObjectInitializer
    ArrayInitializer(mutable bool, elements []Element) ArrayInitializer
    NamedMemberInitializer(name Name, typ Element, value Element) NamedMemberInitializer
    SplatMemberInitializer(typ Element) SplatMemberInitializer
    Lambda(typeParameters TypeParameters, parameters []Parameter, body Element) Lambda
    TypeParameters(parameters []TypeParameter, wheres []Where) TypeParameters
    TypeParameter(name Name, constraint Element) TypeParameter
    Where(left, right Element) Where
    Parameter(name Name, typ Element, deflt Element) Parameter
    VarDefinition(name Name, typ Element, mutable bool) VarDefinition
    Error(message string) Error
    Clone(context BuilderContext) Builder
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
    locations []token.Pos
}

// NewBuilder makes an AST builder that can be used to make AST nodes
func NewBuilder(context BuilderContext) Builder {
    return &builderImpl{context: context}
}

func (b *builderImpl) PushContext() {
    b.locations = append(b.locations, b.context.Start())
}

func (b *builderImpl) PopContext() {
    b.locations = b.locations[0:len(b.locations)-1]
}

func (b *builderImpl) Loc() Location {
    start := b.locations[len(b.locations)-1]
    return NewLocation(start, b.context.End())
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

type literalRuneImpl struct {
    Location
    value rune
}

func (l *literalRuneImpl) Value() rune {
    return l.value
}

func (b *builderImpl) LiteralRune(value rune) LiteralRune {
    return &literalRuneImpl{Location: b.Loc(), value: value}
}

type literalByteImpl struct {
    Location
    value byte
}

func (l *literalByteImpl) Value() byte {
    return l.value
}

func (b *builderImpl) LiteralByte(value byte) LiteralByte {
    return &literalByteImpl{Location: b.Loc(), value: value}
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

type literalUIntImpl struct {
    Location
    value uint
}

func (l *literalUIntImpl) Value() uint {
    return l.value
}

func (b *builderImpl) LiteralUInt(value uint) LiteralUInt {
    return &literalUIntImpl{Location: b.Loc(), value: value}
}

type literalLongImpl struct {
    Location
    value int64
}

func (l *literalLongImpl) Value() int64 {
    return l.value
}

func (b *builderImpl) LiteralLong(value int64) LiteralLong {
    return &literalLongImpl{Location: b.Loc(), value: value}
}

type literalULongImpl struct {
    Location
    value uint64
}

func (l *literalULongImpl) Value() uint64 {
    return l.value
}

func (b *builderImpl) LiteralULong(value uint64) LiteralULong {
    return &literalULongImpl{Location: b.Loc(), value: value}
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

type selectionImpl struct {
    Location
    target Element
    member Name
}

func (l *selectionImpl) Target() Element {
    return l.target
}

func (l *selectionImpl) Member() Name {
    return l.member
}

func (b *builderImpl) Selection(target Element, member Name) Selection {
    return &selectionImpl{Location: b.Loc(), target: target, member: member}
}

type callImpl struct {
    Location
    target Element
    arguments []Element
}

func (c *callImpl) Target() Element {
    return c.target
}

func (c *callImpl) Arguments() []Element {
    return c.arguments
}

func (b *builderImpl) Call(target Element, arguments []Element) Call {
    return &callImpl{Location: b.Loc(), target: target, arguments: arguments}
}

type namedArgumentImpl struct {
    Location
    name Name
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

func (b *builderImpl) NamedArgument(name Name, value Element) NamedArgument {
    return &namedArgumentImpl{Location: b.Loc(), name: name, value: value}
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
    typeParameters TypeParameters
    parameters []Parameter
    body Element
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

func (b *builderImpl) Lambda(typeParameters TypeParameters, parameters []Parameter, body Element) Lambda {
    return &lambdaImpl{Location: b.Loc(), typeParameters: typeParameters, parameters: parameters, body: body}
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

type typeParametersImpl struct {
    Location
    parameters []TypeParameter
    wheres []Where
}

func (t *typeParametersImpl) Parameters() []TypeParameter {
    return t.parameters
}

func (t *typeParametersImpl) Wheres() []Where {
    return t.wheres
}

func (b *builderImpl) TypeParameters(parameters []TypeParameter, wheres []Where) TypeParameters {
    return &typeParametersImpl{Location: b.Loc(), parameters: parameters, wheres: wheres}
}

type typeParameterImpl struct {
    Location
    name Name
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

func (b *builderImpl) TypeParameter(name Name, constraint Element) TypeParameter {
    return &typeParameterImpl{Location: b.Loc(), name: name, constraint: constraint}
}

type whereImpl struct {
    Location
    left Element
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

func (b *builderImpl) Where(left, right Element) Where {
    return &whereImpl{Location: b.Loc(), left: left, right: right}
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

func (b *builderImpl) Clone(context BuilderContext) Builder {
    return &builderImpl{context: context, locations: b.locations}
}

