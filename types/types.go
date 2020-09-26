package types

import (
	"fmt"
	"strings"

	"dyego0/assert"
	"dyego0/symbols"
)

// Type represents the operations that can be
type Type interface {
	// Symbol is the types unique symbol. Other TypeSymbol's might refer to this type but
	// this symbol can be used as the canonical symbol.
	Symbol() TypeSymbol

	// DisplayName is a name that can be used in error messages. DisplayName is allowed to
	// calculate a name and take time doing so, which, for example, would happen if the
	// type is from an anonymous literal. For a quick, non-calculated name, use the symbol's
	// name.
	DisplayName() string

	// Members is an array of the type's members.
	Members() []Member

	// Signatures is an array of callable signature
	Signatures() []Signature
}

// TypeSymbol is a type
type TypeSymbol interface {
	symbols.Symbol

	// The type for which this is the symbol for. This is not necessarily the canoical symbol.
	// The canonical symbol can be retrieved by calling Canonical.
	Type() Type

	// Return the canonical symbol for the type.
	Canonical() TypeSymbol

	// IsType returns true
	IsType() bool
}

// Member is a symbol for the member of a type
type Member interface {
	symbols.Symbol

	// The type of the member
	Type() Type

	// IsMember returns true
	IsMember() bool
}

// Field is a field of a data type
type Field interface {
	Member

	// IsField returns true
	IsField() bool
}

// Signature is a description of the call supported
type Signature interface {
	// This is the context the function is executed in
	This() Type

	// Parameters is the list of parameters for the signature
	Parameters() []Parameter

	// Result is the type of the function result
	Result() Type
}

// Parameter is a function parameter
type Parameter interface {
	symbols.Symbol

	// Type is the type of the parameter
	Type() Type

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
func NewType(symbol TypeSymbol, members []Member, signatures []Signature) Type {
	result := &typeImpl{symbol: symbol, members: members, signatures: signatures}
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

// NewField creates a new field symbol
func NewField(name string, typ Type) Field {
	return &fieldImpl{memberImpl: memberImpl{name: name, typ: typ}}
}

// NewSignature creates a new signature
func NewSignature(this Type, parameters []Parameter, result Type) Signature {
	return &signatureImpl{this: this, parameters: parameters, result: result}
}

// NewParameter creates a new parameter symbol
func NewParameter(name string, typ Type) Parameter {
	return &parameterImpl{name: name, typ: typ}
}

// NewTypeMember creates a new type member
func NewTypeMember(name string, typ Type) TypeMember {
	return &typeMemberImpl{memberImpl: memberImpl{name: name, typ: typ}}
}

type typeImpl struct {
	symbol     TypeSymbol
	members    []Member
	signatures []Signature
}

func (t *typeImpl) Symbol() TypeSymbol {
	return t.symbol
}

func (t *typeImpl) Members() []Member {
	return t.members
}

func (t *typeImpl) Signatures() []Signature {
	return t.signatures
}

func (t *typeImpl) DisplayName() string {
	symbol := t.symbol
	if symbol != nil {
		name := symbol.Name()
		if name != "" {
			return name
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

type memberImpl struct {
	name string
	typ  Type
}

func (m memberImpl) Name() string {
	return m.name
}

func (m memberImpl) Type() Type {
	return m.typ
}

func (m memberImpl) IsMember() bool {
	return true
}

type fieldImpl struct {
	memberImpl
}

func (f *fieldImpl) IsField() bool {
	return true
}

type signatureImpl struct {
	this       Type
	parameters []Parameter
	result     Type
}

func (s *signatureImpl) This() Type {
	return s.this
}

func (s *signatureImpl) Parameters() []Parameter {
	return s.parameters
}

func (s *signatureImpl) Result() Type {
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
	typ  Type
}

func (p *parameterImpl) Name() string {
	return p.name
}

func (p *parameterImpl) Type() Type {
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
