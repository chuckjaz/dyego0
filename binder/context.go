package binder

import (
	"dyego0/errors"
	"dyego0/location"
	"dyego0/symbols"
)

// BindingContext is a context for binding symbols
type BindingContext struct {
	// Scope is the root scope of the context
	Scope symbols.ScopeBuilder

	// Builders is a map of scope symbols to their builders
	Builders map[symbols.Symbol]symbols.ScopeBuilder

	// Errors is the errors reported during binding
	Errors []errors.Error
}

// NewContext creates a new binding context
func NewContext() *BindingContext {
	return &BindingContext{
		Scope:    symbols.NewBuilder(),
		Builders: make(map[symbols.Symbol]symbols.ScopeBuilder),
	}
}

func (context *BindingContext) Error(loc location.Locatable, message string, args ...interface{}) {
	context.Errors = append(context.Errors, errors.New(loc, message, args...))
}
