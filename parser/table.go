package parser

import (
	"fmt"

	"dyego0/assert"
	"dyego0/ast"
)

const infixTypeMember = "infix type member"

type any interface{}

// precedenceLevel

type precedenceLevelImpl struct {
	higher *precedenceLevelImpl
	lower  *precedenceLevelImpl
}

type precedenceLevel interface {
	MakeHigher() precedenceLevel
	MakeLower() precedenceLevel
	Higher() precedenceLevel
	Lower() precedenceLevel
	IsHigherThan(other precedenceLevel) bool
	Level() int
}

func newPrecedenceLevel() precedenceLevel {
	return &precedenceLevelImpl{}
}

func (p *precedenceLevelImpl) Higher() precedenceLevel {
	if p.higher == nil {
		return nil
	}
	return p.higher
}

func (p *precedenceLevelImpl) Lower() precedenceLevel {
	if p.lower == nil {
		return nil
	}
	return p.lower
}

func (p *precedenceLevelImpl) MakeHigher() precedenceLevel {
	newLevel := &precedenceLevelImpl{}
	higher := p.higher
	newLevel.higher = higher
	newLevel.lower = p
	if higher != nil {
		higher.lower = newLevel
	}
	p.higher = newLevel
	return newLevel
}

func (p *precedenceLevelImpl) MakeLower() precedenceLevel {
	newLevel := &precedenceLevelImpl{}
	lower := p.lower
	newLevel.lower = lower
	newLevel.higher = p
	if lower != nil {
		lower.higher = newLevel
	}
	p.lower = newLevel
	return newLevel
}

func (p *precedenceLevelImpl) IsHigherThan(other precedenceLevel) bool {
	var left precedenceLevel = p
	right := other
	if p == other {
		return false
	}
	for left != nil || right != nil {
		if left == other {
			return false
		}
		if right == p {
			return true
		}
		if left != nil {
			left = left.Higher()
		}
		if right != nil {
			right = right.Higher()
		}
	}
	return false
}

func (p *precedenceLevelImpl) Level() int {
	if p.lower != nil {
		return p.lower.Level() + 1
	}
	return 0
}

// operator

type operator interface {
	Name() string
	Levels() []precedenceLevel
	Associativities() []ast.OperatorAssociativity
}

type operatorImpl struct {
	name            string
	levels          []precedenceLevel
	associativities []ast.OperatorAssociativity
}

func (o *operatorImpl) Name() string {
	return o.name
}

func (o *operatorImpl) Levels() []precedenceLevel {
	return o.levels
}

func (o *operatorImpl) Associativities() []ast.OperatorAssociativity {
	return o.associativities
}

func (o *operatorImpl) String() string {
	result := "operator " + o.name
	addPlace := func(placement ast.OperatorPlacement) {
		level := o.levels[placement]
		if level != nil {
			result += fmt.Sprintf(" %s(%d)", placement, level.Level())
		}
	}
	addPlace(ast.Infix)
	addPlace(ast.Prefix)
	addPlace(ast.Postfix)
	return result
}

func newOperator(
	name string,
	levels []precedenceLevel,
	associativities []ast.OperatorAssociativity,
) operator {
	return &operatorImpl{name: name, levels: levels, associativities: associativities}
}

// vocabulary

type vocabularyMap map[string]any

type vocabularyImpl struct {
	members vocabularyMap
	scope   vocabularyScope
}

type vocabulary interface {
	Get(name string) (any, bool)
	Scope() vocabularyScope
}

func newVocabulary() *vocabularyImpl {
	return &vocabularyImpl{members: make(vocabularyMap), scope: newVocabularyScope()}
}

func (v *vocabularyImpl) Get(name string) (any, bool) {
	result, ok := v.members[name]
	return result, ok
}

func (v *vocabularyImpl) Scope() vocabularyScope {
	return v.scope
}

func (v *vocabularyImpl) String() string {
	result := "vocabulary:\n"
	for _, element := range v.members {
		result += fmt.Sprintf("  %s\n", element)
	}
	return result
}

// vocabularyScope

type vocabularyScopeImpl struct {
	members map[string]any
}

type vocabularyScope interface {
	Get(name string) (any, bool)
}

func newVocabularyScope() *vocabularyScopeImpl {
	return &vocabularyScopeImpl{members: make(map[string]any)}
}

func (v *vocabularyScopeImpl) Get(name string) (any, bool) {
	value, ok := v.members[name]
	return value, ok
}

// vocabularyError

type vocabularyError struct {
	element ast.Element
	message string
}

type vocabularyErrors []vocabularyError

// vocabularyEmbeddingContext
type vocabularyEmbeddingContext struct {
	result        *vocabularyImpl
	errors        vocabularyErrors
	rootLevel     precedenceLevel
	lowestLevel   precedenceLevel
	precedenceMap map[precedenceLevel]precedenceLevel
}

func newVocabularyEmbeddingContext() *vocabularyEmbeddingContext {
	level := newPrecedenceLevel()
	return &vocabularyEmbeddingContext{
		result:        newVocabulary(),
		rootLevel:     level,
		lowestLevel:   level,
		precedenceMap: make(map[precedenceLevel]precedenceLevel),
	}
}

func (c *vocabularyEmbeddingContext) mappedPrecedence(precedence precedenceLevel) precedenceLevel {
	if precedence == nil {
		return c.rootLevel
	}
	level, ok := c.precedenceMap[precedence]
	if ok {
		return level
	}
	parent := c.mappedPrecedence(precedence.Higher())
	result := parent.MakeLower()
	c.precedenceMap[precedence] = result
	return result
}

func (c *vocabularyEmbeddingContext) mappedPrecedences(precedence []precedenceLevel) []precedenceLevel {
	result := []precedenceLevel{nil, nil, nil}
	for placement, level := range precedence {
		if level != nil {
			result[placement] = c.mappedPrecedence(level)
		}
	}
	return result
}

func (c *vocabularyEmbeddingContext) reportError(element ast.Element, message string, args ...interface{}) {
	c.errors = append(c.errors, vocabularyError{
		element: element,
		message: fmt.Sprintf(message, args...),
	})
}

func (c *vocabularyEmbeddingContext) embedVocabulary(embeddedVocabulary *vocabularyImpl, embedding ast.Element) {
	members := embeddedVocabulary.members
	for _, member := range members {
		switch m := member.(type) {
		case operator:
			c.recordOperator(embedding, m.Name(), c.mappedPrecedences(m.Levels()), m.Associativities())
		}
	}
	last := c.rootLevel
	for {
		next := last.Lower()
		if next == nil {
			break
		}
		last = next
	}
	c.lowestLevel = last
}

func (c *vocabularyEmbeddingContext) recordOperator(
	element ast.Element,
	name string,
	levels []precedenceLevel,
	associativities []ast.OperatorAssociativity,
) {
	member, ok := c.result.Get(name)
	if ok {
		switch m := member.(type) {
		case operator:
			for placement := ast.OperatorPlacement(0); placement < ast.UnspecifiedPlacement; placement++ {
				if levels[placement] != nil {
					if m.Levels()[placement] == nil {
						m.Levels()[placement] = levels[placement]
						m.Associativities()[placement] = associativities[placement]
					} else {
						c.reportError(
							element,
							"An %s operator '%s' already defined",
							placement.String(),
							m.Name(),
						)
					}
				}
			}
		}
		return
	}
	op := newOperator(
		name,
		levels,
		associativities,
	)
	c.result.members[name] = op
}

// buildVocabulary

func buildVocabulary(scope vocabularyScope, vocabularyLiteral ast.VocabularyLiteral) (vocabulary, vocabularyErrors) {
	c := newVocabularyEmbeddingContext()

	lookupVocabulary := func(nameList []ast.Name) vocabulary {
		currentScope := scope
		var embeddedVocabulary vocabulary
		var lastName ast.Name
		for _, name := range nameList {
			if currentScope == nil {
				c.reportError(lastName, "Expected '%s' to be a vocabulary scope", lastName.Text())
				return nil
			}
			lookup, ok := currentScope.Get(name.Text())
			if !ok {
				c.reportError(name, "Undefined vocabulary '%s'", name.Text())
				return nil
			}
			switch v := lookup.(type) {
			case vocabulary:
				embeddedVocabulary = v
				currentScope = nil
			case vocabularyScope:
				currentScope = v
			default:
				assert.Fail("Unknown scope member %#v", lookup)
			}
			lastName = name
		}
		if embeddedVocabulary == nil {
			c.reportError(lastName, "Expected '%s' to be a vocabulary", lastName.Text())
			return nil
		}
		return embeddedVocabulary
	}

	findLowestPrecedence := func() precedenceLevel {
		current := c.rootLevel
		last := c.rootLevel
		for current != nil {
			last = current
			current = current.Lower()
		}
		return last
	}

	// Resolve any embedded vocabularies
	members := vocabularyLiteral.Members()
	for _, member := range members {
		switch m := member.(type) {
		case ast.VocabularyEmbedding:
			embeddedVocabulary := lookupVocabulary(m.Name())
			if embeddedVocabulary == nil {
				continue
			}
			c.embedVocabulary(embeddedVocabulary.(*vocabularyImpl), m)
		case ast.VocabularyOperatorDeclaration:
			continue
		default:
			assert.Fail("Unknown vmocabulary element %#v", m)
		}
	}

	levelsAndAssociativities := func(
		placement ast.OperatorPlacement,
		precedence precedenceLevel,
		associativity ast.OperatorAssociativity,
	) ([]precedenceLevel, []ast.OperatorAssociativity) {
		levels := []precedenceLevel{nil, nil, nil}
		levels[placement] = precedence
		associativities := []ast.OperatorAssociativity{
			ast.UnspecifiedAssociativity,
			ast.UnspecifiedAssociativity,
			ast.UnspecifiedAssociativity,
		}
		associativities[placement] = associativity
		return levels, associativities
	}

	// Declare operators
	lowestPrecedence := findLowestPrecedence()
	for _, member := range members {
		switch m := member.(type) {
		case ast.VocabularyEmbedding:
			continue
		case ast.VocabularyOperatorDeclaration:
			placement := m.Placement()
			associativity := m.Associativity()
			precedence := lowestPrecedence
			precedenceDeclaration := m.Precedence()
			if precedenceDeclaration != nil {
				lookup, ok := c.result.members[precedenceDeclaration.Name().Text()]
				if !ok {
					c.reportError(precedenceDeclaration.Name(), "Undeclared identifier '%s'", precedenceDeclaration.Name().Text())
					continue
				}
				referencedOperator, ok := lookup.(operator)
				if !ok {
					c.reportError(precedenceDeclaration.Name(), "'%s' does not refer to an operator", precedenceDeclaration.Name().Text())
					continue
				}
				referencedPlacement := precedenceDeclaration.Placement()
				if referencedPlacement == ast.UnspecifiedPlacement {
					for placement := ast.OperatorPlacement(0); placement < ast.UnspecifiedPlacement; placement++ {
						if referencedOperator.Associativities()[placement] != ast.UnspecifiedAssociativity {
							if referencedPlacement != ast.UnspecifiedPlacement {
								c.reportError(
									precedenceDeclaration,
									"Ambigious operator reference, both %s and %s are defined",
									referencedPlacement,
									placement,
								)
							}
							referencedPlacement = placement
						}
					}
				}
				precedence = referencedOperator.Levels()[referencedPlacement]
				if precedence == nil {
					c.reportError(
						precedenceDeclaration,
						"No %s placement defined for operator '%s'",
						referencedPlacement,
						precedenceDeclaration.Name().Text(),
					)
					continue
				}
				switch precedenceDeclaration.Relation() {
				case ast.Before:
					precedence = precedence.MakeHigher()
				case ast.After:
					precedence = precedence.MakeLower()
				default:
					assert.Fail("Relation not defined: %s", precedenceDeclaration.Relation())
				}
			} else {
				precedence := lowestPrecedence.MakeLower()
				lowestPrecedence = precedence
			}

			for _, name := range m.Names() {
				if name.Text() == infixTypeMember && placement != ast.Infix {
					c.reportError(
						precedenceDeclaration,
						"Local type identifiers can only be declaried as an infix operator",
					)
					continue
				}
				levels, associativities := levelsAndAssociativities(placement, precedence, associativity)
				c.recordOperator(name, name.Text(), levels, associativities)
			}
		}
	}
	return c.result, c.errors
}
