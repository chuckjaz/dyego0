package types

import (
	"fmt"
	"strings"

	"dyego0/assert"
	"dyego0/symbols"
)

// TypeKind is the kind of type
type TypeKind int

const (
	// Record is a linear block of memory separated into files
	Record TypeKind = iota

	// Reference to a record or array
	Reference

	// Array is a linear block of memory of homomorphic type
	Array

	// Module is a fixed block of memory similar to a record but in static memory
	Module

	// Error is the type of invalid expressions
	Error
)

// Type represents the operations that can be
type Type interface {
	// Symbol is the types unique symbol. Other TypeSymbol's might refer to this type but
	// this symbol can be used as the canonical symbol.
	Symbol() TypeSymbol

	// Kind is the kind of the type
	Kind() TypeKind

	// DisplayName is a name that can be used in error messages. DisplayName is allowed to
	// calculate a name and take time doing so, which, for example, would happen if the
	// type is from an anonymous literal. For a quick, non-calculated name, use the symbol's
	// name.
	DisplayName() string

	// Members is an array of the type's instance members.
	Members() []Member

	// MemberScope is a lookup scope for instqnce members
	MemberScope() symbols.Scope

	// TypeScope is the lookup scope for type members
	TypeScope() symbols.Scope

	// Signatures is an array of callable signature
	Signatures() []Signature

	// Container is the type that contains this type if there is one
	Container() TypeSymbol

	// Elements is the element type of an Array
	Elements() TypeSymbol

	// Size is the size of the array
	Size() int

	// Referant is the type a reference refers to
	Referant() TypeSymbol

	// String returns the display name
	String() string
}

// TypeSymbol is a type
// A type symbol with a nil type can act as an open type variable and be bound using UpdateTypeSymbol
// If a Type is created from it using NewType, this happens during NewType
type TypeSymbol interface {
	symbols.Symbol

	// The type for which this is the symbol for. This is not necessarily the canoical symbol.
	// The canonical symbol can be retrieved by calling Canonical.
	Type() Type

	// Return the canonical symbol for the type.
	Canonical() TypeSymbol

	// IsType returns true
	IsType() bool

	// String is the display name of the type symbol
	String() string
}

// Member is a symbol for the member of a type
type Member interface {
	symbols.Symbol

	// The type of the member
	Type() TypeSymbol

	// IsMember returns true
	IsMember() bool
}

// Field is a field of a data type
type Field interface {
	Member

	// Mutable is true if the field can be mutated
	Mutable() bool

	// IsField returns true
	IsField() bool
}

// Signature is a description of the call supported
type Signature interface {
	// This is the context the function is executed in
	This() TypeSymbol

	// Parameters is the list of parameters for the signature
	Parameters() []Parameter

	// Result is the type of the function result
	Result() TypeSymbol
}

// Parameter is a function parameter
type Parameter interface {
	symbols.Symbol

	// Type is the type of the parameter
	Type() TypeSymbol

	// IsParameter returns true
	IsParameter() bool
}

// TypeMember is an embedded type
type TypeMember interface {
	Member

	// IsTypeMember returns true
	IsTypeMember() bool
}

// NewType create a new type
func NewType(
	symbol TypeSymbol,
	kind TypeKind,
	members []Member,
	memberScope symbols.Scope,
	typeScope symbols.Scope,
	signatures []Signature,
	container TypeSymbol,
) Type {
	if members != nil && memberScope == nil {
		b := symbols.NewBuilder()
		for _, member := range members {
			b.Enter(member)
		}
		memberScope = b.Build()
	}
	if memberScope == nil {
		memberScope = symbols.EmptyScope()
	}
	if typeScope == nil {
		typeScope = symbols.EmptyScope()
	}
	result := &typeImpl{
		symbol:      symbol,
		kind:        kind,
		members:     members,
		memberScope: memberScope,
		typeScope:   typeScope,
		signatures:  signatures,
		container:   container,
	}
	if symbol.Type() == nil {
		UpdateTypeSymbol(symbol, result)
	}
	return result
}

// NewTypeSymbol creates a new type symbol. Passing nil for typ is allowed when creating a named type
func NewTypeSymbol(name string, typ Type) TypeSymbol {
	return &typeSymbolImpl{name: name, typ: typ}
}

// UpdateTypeSymbol can update a type symbol with a nil type
func UpdateTypeSymbol(sym TypeSymbol, typ Type) {
	ts := sym.(*typeSymbolImpl)
	assert.Assert(ts.typ == nil, "Trying to update an aready defined type")
	ts.typ = typ
}

// NewArrayType creates a new array type
func NewArrayType(symbol TypeSymbol, elements TypeSymbol, size int) Type {
	result := &arrayType{
		symbol:   symbol,
		elements: elements,
		size:     size,
	}
	if symbol.Type() == nil {
		UpdateTypeSymbol(symbol, result)
	}
	return result
}

// MakeArray creates a new array
func MakeArray(elements TypeSymbol) TypeSymbol {
	result := NewTypeSymbol(fmt.Sprintf("%s[]", elements.Name()), nil)
	NewArrayType(result, elements, -1)
	return result
}

// NewReferenceType creates a new reference type
func NewReferenceType(symbol TypeSymbol, referant TypeSymbol) Type {
	result := &referenceType{
		symbol:   symbol,
		referant: referant,
	}
	if symbol.Type() == nil {
		UpdateTypeSymbol(symbol, result)
	}
	return result
}

// MakeReference makes a reference type that references referant
func MakeReference(referant TypeSymbol) TypeSymbol {
	r := NewTypeSymbol("*"+referant.Name(), nil)
	NewReferenceType(r, referant)
	return r
}

// NewField creates a new field symbol
func NewField(name string, typ TypeSymbol, mutable bool) Field {
	return &fieldImpl{memberImpl: memberImpl{name: name, typ: typ}, mutable: mutable}
}

// NewSignature creates a new signature
func NewSignature(this TypeSymbol, parameters []Parameter, result TypeSymbol) Signature {
	return &signatureImpl{this: this, parameters: parameters, result: result}
}

// NewParameter creates a new parameter symbol
func NewParameter(name string, typ TypeSymbol) Parameter {
	return &parameterImpl{name: name, typ: typ}
}

// NewTypeMember creates a new type member
func NewTypeMember(name string, typ TypeSymbol) TypeMember {
	return &typeMemberImpl{memberImpl: memberImpl{name: name, typ: typ}}
}

// NewErrorType creates an error type symbol
func NewErrorType() TypeSymbol {
	result := NewTypeSymbol("<error>", nil)
	NewType(result, Error, nil, nil, nil, nil, nil)
	return result
}

// IsError returns true if typeSym is an error type
func IsError(typeSym TypeSymbol) bool {
	t := typeSym.Type()
	return t != nil && t.Kind() == Error
}

type typeImpl struct {
	symbol      TypeSymbol
	kind        TypeKind
	members     []Member
	memberScope symbols.Scope
	typeScope   symbols.Scope
	signatures  []Signature
	container   TypeSymbol
}

func (t *typeImpl) Symbol() TypeSymbol {
	return t.symbol
}

func (t *typeImpl) Kind() TypeKind {
	return t.kind
}

func (t *typeImpl) Members() []Member {
	return t.members
}

func (t *typeImpl) MemberScope() symbols.Scope {
	return t.memberScope
}

func (t *typeImpl) TypeScope() symbols.Scope {
	return t.typeScope
}

func (t *typeImpl) Signatures() []Signature {
	return t.signatures
}

func (t *typeImpl) Container() TypeSymbol {
	return t.container
}

func (t *typeImpl) containerName() string {
	if t.container == nil {
		return ""
	}
	ct := t.container.Type()
	if ct != nil {
		return ct.DisplayName() + "."
	}
	n := t.container.Name()
	if n != "" {
		return n + "."
	}
	return ""
}

func (t *typeImpl) DisplayName() string {
	symbol := t.symbol
	if symbol != nil {
		name := symbol.Name()
		if name != "" {
			return t.containerName() + name
		}
	}
	builder := &stringBuilder{}
	builder.List("<", ">", func() {
		for _, member := range t.members {
			builder.Item(func() {
				builder.Add(member.Name())
				builder.Add(": ")
				builder.Convert(member.Type())
			})
		}
		for _, signature := range t.signatures {
			builder.Item(func() {
				builder.Convert(signature)
			})
		}
	})
	return builder.String()
}

func (t *typeImpl) Elements() TypeSymbol {
	return nil
}

func (t *typeImpl) Size() int {
	return 0
}

func (t *typeImpl) Referant() TypeSymbol {
	return nil
}

func (t *typeImpl) String() string {
	return t.DisplayName()
}

type typeSymbolImpl struct {
	name string
	typ  Type
}

func (s *typeSymbolImpl) Name() string {
	return s.name
}

func (s *typeSymbolImpl) Type() Type {
	return s.typ
}

func (s *typeSymbolImpl) IsType() bool {
	return true
}

func (s *typeSymbolImpl) Canonical() TypeSymbol {
	return s.Type().Symbol()
}

func (s *typeSymbolImpl) String() string {
	if s.typ != nil {
		return s.typ.DisplayName()
	}
	return s.name
}

type arrayType struct {
	symbol   TypeSymbol
	elements TypeSymbol
	size     int
}

func (a *arrayType) Symbol() TypeSymbol {
	return a.symbol
}

func (a *arrayType) Kind() TypeKind {
	return Array
}

func (a *arrayType) DisplayName() string {
	elementName := a.elements.String()
	if a.size < 0 {
		return elementName + "[]"
	}
	return fmt.Sprintf("%s[%d]", elementName, a.size)
}

func (a *arrayType) Members() []Member {
	return nil
}

func (a *arrayType) MemberScope() symbols.Scope {
	return symbols.EmptyScope()
}

func (a *arrayType) TypeScope() symbols.Scope {
	return symbols.EmptyScope()
}

func (a *arrayType) Signatures() []Signature {
	return nil
}

func (a *arrayType) Extensions() symbols.Scope {
	return symbols.EmptyScope()
}

func (a *arrayType) Container() TypeSymbol {
	return nil
}

func (a *arrayType) Elements() TypeSymbol {
	return a.elements
}

func (a *arrayType) Size() int {
	return a.size
}

func (a *arrayType) Referant() TypeSymbol {
	return nil
}

func (a *arrayType) String() string {
	return a.DisplayName()
}

type referenceType struct {
	symbol   TypeSymbol
	referant TypeSymbol
}

func (r *referenceType) Symbol() TypeSymbol {
	return r.symbol
}

func (r *referenceType) Kind() TypeKind {
	return Reference
}

func (r *referenceType) DisplayName() string {
	return "*" + r.referant.String()
}

func (r *referenceType) Members() []Member {
	return nil
}

func (r *referenceType) MemberScope() symbols.Scope {
	return symbols.EmptyScope()
}

func (r *referenceType) TypeScope() symbols.Scope {
	return symbols.EmptyScope()
}

func (r *referenceType) Signatures() []Signature {
	return nil
}

func (r *referenceType) Extensions() symbols.Scope {
	return symbols.EmptyScope()
}

func (r *referenceType) Container() TypeSymbol {
	return nil
}

func (r *referenceType) Elements() TypeSymbol {
	return nil
}

func (r *referenceType) Size() int {
	return 0
}

func (r *referenceType) Referant() TypeSymbol {
	return r.referant
}

func (r *referenceType) String() string {
	return r.DisplayName()
}

type memberImpl struct {
	name string
	typ  TypeSymbol
}

func (m memberImpl) Name() string {
	return m.name
}

func (m memberImpl) Type() TypeSymbol {
	return m.typ
}

func (m memberImpl) IsMember() bool {
	return true
}

type fieldImpl struct {
	memberImpl
	mutable bool
}

func (f *fieldImpl) Mutable() bool {
	return f.mutable
}

func (f *fieldImpl) IsField() bool {
	return true
}

type signatureImpl struct {
	this       TypeSymbol
	parameters []Parameter
	result     TypeSymbol
}

func (s *signatureImpl) This() TypeSymbol {
	return s.this
}

func (s *signatureImpl) Parameters() []Parameter {
	return s.parameters
}

func (s *signatureImpl) Result() TypeSymbol {
	return s.result
}

func (s *signatureImpl) String() string {
	builder := &stringBuilder{}
	builder.List("{", "}", func() {
		if s.this != nil {
			builder.Convert(s.this)
			builder.Add(".")
		}
		for _, parameter := range s.parameters {
			builder.Item(func() {
				builder.Add(parameter.Name())
				builder.Add(": ")
				builder.Convert(parameter.Type())
			})
		}
		builder.Add(" -> ")
		builder.Convert(s.Result())
	})
	return builder.String()
}

type parameterImpl struct {
	name string
	typ  TypeSymbol
}

func (p *parameterImpl) Name() string {
	return p.name
}

func (p *parameterImpl) Type() TypeSymbol {
	return p.typ
}

func (p *parameterImpl) IsParameter() bool {
	return true
}

type typeMemberImpl struct {
	memberImpl
}

func (t *typeMemberImpl) IsTypeMember() bool {
	return true
}

type stringBuilder struct {
	result []string
	first  bool
}

func (s *stringBuilder) Add(value string) {
	s.result = append(s.result, value)
}

func (s *stringBuilder) Convert(value interface{}) {
	s.Add(fmt.Sprintf("%s", value))
}

func (s *stringBuilder) String() string {
	return strings.Join(s.result, "")
}

func (s *stringBuilder) List(start, end string, block func()) {
	s.first = true
	s.Add(start)
	block()
	s.Add(end)
}

func (s *stringBuilder) Item(block func()) {
	if !s.first {
		s.Add(", ")
	}
	s.first = false
	block()
}
