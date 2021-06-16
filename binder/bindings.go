package binder

import (
	"dyego0/assert"
	"dyego0/types"
)

// Bindings is a unification context
type Bindings interface {
	UnifyType(t, v types.TypeSymbol) bool
}

type binding struct {
	typ  types.TypeSymbol
	size *int
	next *binding
}

// NewBindings creates a new unification context
func NewBindings() Bindings {
	return &bindings{
		bindings: make(map[types.TypeSymbol]*binding),
	}
}

type bindings struct {
	bindings map[types.TypeSymbol]*binding
}

func (b *bindings) remember(typSym types.TypeSymbol) *binding {
	v, ok := b.bindings[typSym]
	if !ok {
		size := 1
		v = &binding{typ: typSym, size: &size}
		v.next = v
		b.bindings[typSym] = v
	}
	return v
}

func (b *bindings) bindType(typSym types.TypeSymbol, typ types.Type) {
	assert.Assert(typSym.Type() == nil, "Cannot bind a bound type")
	types.UpdateTypeSymbol(typSym, typ)

	// If the type was bound to other variables, update them too
	v, ok := b.bindings[typSym]
	if ok {
		cur := v.next
		for cur != v {
			types.UpdateTypeSymbol(cur.typ, typ)
			cur = cur.next
		}
	}
}

func updateSize(b *binding, size *int) {
	b.size = size
	cur := b.next
	for cur != b {
		cur.size = size
		cur = cur.next
	}
}

func (b *bindings) bind(t, v types.TypeSymbol) bool {
	tv := b.remember(t)
	vv := b.remember(v)

	// If the size references are the same then the symbols are already bound together
	if tv.size == vv.size {
		return true
	}

	// Update the smaller list to use largers size
	if *tv.size < *vv.size {
		*tv.size += *vv.size
		updateSize(vv, tv.size)
	} else {
		*vv.size += *tv.size
		updateSize(tv, vv.size)
	}

	// Merge the lists by swaping their next pointers
	tvn := tv.next
	vvn := vv.next
	tv.next = vvn
	vv.next = tvn

	return true
}

func (b *bindings) UnifyType(t, v types.TypeSymbol) bool {
	tt := t.Type()
	vt := v.Type()
	if tt == nil && vt != nil {
		// type t is open, v is closed update the types
		b.bindType(t, vt)
		return true
	} else if tt != nil && vt == nil {
		// type v is open and t is closed
		b.bindType(v, tt)
		return true
	} else if tt != nil && vt != nil {
		if tt == vt {
			return true
		}
		tk := tt.Kind()
		vk := vt.Kind()
		if tk == vk {
			switch tt.Kind() {
			case types.Array:
				return b.UnifyType(tt.Elements(), vt.Elements())
			case types.Reference:
				return b.UnifyType(tt.Referant(), vt.Referant())
			}
		}
		return false
	}

	// both are unbound but need be bound together
	return b.bind(v, t)
}
