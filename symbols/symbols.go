package symbols

// Symbol is a named reference to something
type Symbol interface {
	Name() string
}

// Scope is a read-only map of names to symbols
type Scope interface {
	// Find finds a symbol in a scope. Returns the symbol matching name and true or nil and false
	Find(name string) (Symbol, bool)

	// Determine if name is in scope
	Contains(name string) bool

	// ForEach enumerates the symbols in scope
	ForEach(block func(Symbol) bool)

	// Return true if the scope is empty
	IsEmpty() bool
}

// ScopeBuilder is used to build a scope.
type ScopeBuilder interface {
	Scope

	// Enter enters a symbol into scope, Returns the symbol and true if it not already in the table
	// or the prevoius symbol and false, it the symbol already exists.
	Enter(symbol Symbol) (Symbol, bool)

	// Reenter enters the symbol unconditionally into the table overwriting previous symbol with the
	// same name if one is already in the table.
	Reenter(symbol Symbol)

	// Build will return a the scope built with this builder.
	Build() Scope
}

// EmptyScope returns an empty scope
func EmptyScope() Scope {
	return emptyScope{}
}

type scope struct {
	table map[string]Symbol
}

func (s *scope) Find(name string) (Symbol, bool) {
	r, ok := s.table[name]
	return r, ok
}

func (s *scope) Contains(name string) bool {
	_, ok := s.table[name]
	return ok
}

func (s *scope) ForEach(block func(Symbol) bool) {
	for _, symbol := range s.table {
		if block(symbol) {
			break
		}
	}
}

func (s *scope) IsEmpty() bool {
	return len(s.table) == 0
}

func newScope(table map[string]Symbol) Scope {
	return &scope{table: table}
}

type scopeBuilder struct {
	base  Scope
	table map[string]Symbol
}

func (s *scopeBuilder) Enter(symbol Symbol) (Symbol, bool) {
	name := symbol.Name()
	previous, ok := s.table[name]
	if ok {
		return previous, false
	}
	previous, ok = s.base.Find(name)
	if ok {
		return previous, false
	}
	s.table[name] = symbol
	return symbol, true
}

func (s *scopeBuilder) Reenter(symbol Symbol) {
	s.table[symbol.Name()] = symbol
}

func (s *scopeBuilder) Find(name string) (Symbol, bool) {
	r, ok := s.table[name]
	if ok {
		return r, true
	}
	r, ok = s.base.Find(name)
	return r, ok
}

func (s *scopeBuilder) Contains(name string) bool {
	_, ok := s.table[name]
	if ok {
		return true
	}
	ok = s.base.Contains(name)
	return ok
}

func (s *scopeBuilder) ForEach(block func(Symbol) bool) {
	for _, symbol := range s.table {
		if block(symbol) {
			return
		}
	}
	s.base.ForEach(block)
}

func (s *scopeBuilder) IsEmpty() bool {
	if len(s.table) != 0 {
		return false
	}
	return s.base.IsEmpty()
}

func (s *scopeBuilder) Build() Scope {
	table := s.table
	s.table = nil
	_, ok := s.base.(emptyScope)
	if ok {
		return newScope(table)
	}
	return newMultiScope(newScope(table), s.base)
}

// NewBuilderFrom creates a new scope builder that contains all the entries of the given base scope.
func NewBuilderFrom(base Scope) ScopeBuilder {
	table := make(map[string]Symbol)
	s, ok := base.(*scope)
	if ok {
		for k, v := range s.table {
			table[k] = v
		}
		base = emptyScope{}
	}
	return &scopeBuilder{table: table, base: base}
}

// NewBuilder create an empty symbol table builder
func NewBuilder() ScopeBuilder {
	return &scopeBuilder{table: make(map[string]Symbol), base: emptyScope{}}
}

type multiScope struct {
	scopes []Scope
}

func (s *multiScope) Find(name string) (Symbol, bool) {
	for _, scope := range s.scopes {
		result, ok := scope.Find(name)
		if ok {
			return result, true
		}
	}
	return nil, false
}

func (s *multiScope) Contains(name string) bool {
	_, ok := s.Find(name)
	return ok
}

func (s *multiScope) ForEach(block func(Symbol) bool) {
	emitted := make(map[Symbol]bool)
	exit := false
	for _, table := range s.scopes {
		table.ForEach(func(symbol Symbol) bool {
			_, ok := emitted[symbol]
			if !ok {
				emitted[symbol] = true
				exit = block(symbol)
			}
			return exit
		})
		if exit {
			break
		}
	}
}

func (s *multiScope) IsEmpty() bool {
	for _, table := range s.scopes {
		if !table.IsEmpty() {
			return false
		}
	}
	return true
}

func newMultiScope(scopes ...Scope) Scope {
	var result []Scope
	for _, scope := range scopes {
		switch s := scope.(type) {
		case *multiScope:
			result = append(result, s.scopes...)
		case emptyScope:
			// Ignored
		default:
			result = append(result, scope)
		}
	}
	return &multiScope{scopes: result}
}

type emptyScope struct {
}

func (s emptyScope) Find(name string) (Symbol, bool) {
	return nil, false
}

func (s emptyScope) Contains(name string) bool {
	return false
}

func (s emptyScope) ForEach(block func(Symbol) bool) {
}

func (s emptyScope) IsEmpty() bool {
	return true
}

// Merge merges scopes where the earlier scope takes precedence over later scopes
func Merge(scopes ...Scope) Scope {
	switch len(scopes) {
	case 0:
		return emptyScope{}
	case 1:
		return scopes[0]
	default:
		return newMultiScope(scopes...)
	}
}
