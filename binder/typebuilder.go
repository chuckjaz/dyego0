package binder

import (
	"dyego0/symbols"
	"dyego0/types"
)

// typeBuilder

type typeBuilder struct {
	membersScope   symbols.ScopeBuilder
	members        []types.Member
	signatures     []types.Signature
	symbol         types.TypeSymbol
	typeScope      symbols.ScopeBuilder
	nestedBuilders map[symbols.Symbol]*typeBuilder
}

func (b *typeBuilder) AddMember(member types.Member) bool {
	_, ok := b.membersScope.Enter(member)
	b.members = append(b.members, member)
	return ok
}

func (b *typeBuilder) AddSignature(signature types.Signature) {
	b.signatures = append(b.signatures, signature)
}

func (b *typeBuilder) AddTypeSymbol(symbol symbols.Symbol) bool {
	_, ok := b.typeScope.Enter(symbol)
	return ok
}

func (b *typeBuilder) RecordNestedTypeBuilder(symbol symbols.Symbol, nestedBuilder *typeBuilder) {
	b.nestedBuilders[symbol] = nestedBuilder
}

func (b *typeBuilder) FindNestedTypeBuilder(symbol symbols.Symbol) (*typeBuilder, bool) {
	typeBuilder, ok := b.nestedBuilders[symbol]
	return typeBuilder, ok
}

func (b *typeBuilder) FindTypeSymbol(name string) (symbols.Symbol, bool) {
	symbol, ok := b.typeScope.Find(name)
	return symbol, ok
}

func (b *typeBuilder) Build() types.Type {
	membersScope := b.membersScope.Build()
	typeScope := b.typeScope.Build()
	return types.NewType(b.symbol, b.members, membersScope, b.signatures, typeScope)
}

func newTypeBuilder(symbol types.TypeSymbol) *typeBuilder {
	membersScope := symbols.NewBuilder()
	typeScope := symbols.NewBuilder()
	nestedTypeBuilder := make(map[symbols.Symbol]*typeBuilder)
	return &typeBuilder{
		symbol:         symbol,
		membersScope:   membersScope,
		typeScope:      typeScope,
		nestedBuilders: nestedTypeBuilder,
	}
}
