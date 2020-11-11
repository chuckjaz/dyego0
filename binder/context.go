package binder

import (
	"dyego0/errors"
	"dyego0/symbols"
)

// BindingContext is a context for binding symbols
type BindingContext struct {
	// Scope is the root scope of the context
	Scope symbols.ScopeBuilder

	// Errors is the errors reported during binding
	Errors []errors.Error
}

// NewContext creates a new binding context
func NewContext() *BindingContext {
	return &BindingContext{Scope: symbols.NewBuilder()}
}
