package parser

import (
	"fmt"

	"dyego0/assert"
	"dyego0/ast"
	"dyego0/errors"
	"dyego0/location"
	"dyego0/scanner"
	"dyego0/tokens"
)

// Parser parses text and returns an ast element
type Parser interface {
	Errors() []errors.Error
	Parse() ast.Element
}

type parser struct {
	scanner           *scanner.Scanner
	builder           ast.Builder
	current           tokens.Token
	pseudo            tokens.PseudoToken
	operator          *selectedOperator
	excludedOperators []operator
	separatorState    separatorState
	scope             vocabularyScope
	vocabulary        vocabulary
	embeddingContext  *vocabularyEmbeddingContext
	errors            []errors.Error
}

type separatorState int

const (
	normalState separatorState = iota
	wasInfixState
	separatorImplied
)

func (s separatorState) String() string {
	switch s {
	case normalState:
		return "normalState"
	case wasInfixState:
		return "wasInfixState"
	case separatorImplied:
		return "separatorImplied"
	}
	return "!invlidSeparatorState"
}

// NewParser creates a new parser
func NewParser(scanner *scanner.Scanner, scope vocabularyScope) Parser {
	builder := ast.NewBuilder(scanner)
	context := newVocabularyEmbeddingContext()
	p := &parser{
		scanner:          scanner,
		builder:          builder,
		pseudo:           tokens.InvalidPseudoToken,
		separatorState:   normalState,
		scope:            scope,
		vocabulary:       context.result,
		embeddingContext: context,
	}
	builder.PushContext()
	p.next()
	return p
}

func (p *parser) Parse() ast.Element {
	expr := p.sequence()
	p.expect(tokens.EOF)
	return expr
}

func (p *parser) Errors() []errors.Error {
	return p.errors
}

func (p *parser) report(msg string, args ...interface{}) errors.Error {
	err := p.builder.Error(msg, args...)
	errors := p.errors
	l := len(errors)
	if l == 0 || errors[l-1].Start() != err.Start() {
		p.errors = append(p.errors, err)
	}
	return err
}

func (p *parser) reportElement(element ast.Element, msg string, args ...interface{}) ast.Element {
	err := errors.New(element, msg, args...)
	p.errors = append(p.errors, err)
	return err
}

func (p *parser) expect(t tokens.Token) {
	p.builder.PushContext()
	defer p.builder.PopContext()
	if p.current == t {
		p.next()
	} else {
		p.builder.PushContext()
		defer p.builder.PopContext()
		p.report("Expected %v, received %v", t, p.current)
		p.next()
	}
}

func (p *parser) expectPseudo(t tokens.PseudoToken) {
	p.builder.PushContext()
	defer p.builder.PopContext()
	if p.pseudo == t {
		p.next()
	} else {
		p.builder.PushContext()
		defer p.builder.PopContext()
		if (p.current == tokens.Identifier || p.current == tokens.Symbol) && p.pseudo != tokens.InvalidPseudoToken {
			p.report("Expected %s, received %s", t.String(), p.pseudo.String())
		} else {
			p.report("Expected %s, received %v", t.String(), p.current.String())
		}
	}
}

func (p *parser) expects(ts ...tokens.Token) ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	first := true
	result := ""
	for _, t := range ts {
		if !first {
			result += ", "
		}
		result += t.String()
		first = false
	}
	p.report("Expected one of %s, received %v", result, p.current)
	p.next()
	return p.errors[len(p.errors)-1]
}

func (p *parser) expectsPseudo(ts ...tokens.PseudoToken) ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	first := true
	result := ""
	for _, t := range ts {
		if !first {
			result += ", "
		}
		result += t.String()
		first = false
	}
	if p.current == tokens.Identifier && p.pseudo != tokens.InvalidPseudoToken {
		p.report("Expected one of %s, received %s", result, p.pseudo)
	} else {
		p.report("Expected one of %s, received %s", result, p.current)
	}
	p.next()
	return p.errors[len(p.errors)-1]
}

func (p *parser) expectItems(items ...interface{}) ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	first := true
	result := ""
	for _, t := range items {
		if !first {
			result += ", "
		}
		result += fmt.Sprintf("%s", t)
		first = false
	}
	if p.current == tokens.Identifier && p.pseudo != tokens.InvalidPseudoToken {
		p.report("Expected one of %s, received %s", result, p.pseudo)
	} else {
		p.report("Expected one of %s, received %s", result, p.current)
	}
	p.next()
	return p.errors[len(p.errors)-1]
}

func (p *parser) next() tokens.Token {
	var next = p.scanner.Next()
	p.current = next
	p.separatorState = normalState
	p.operator = nil
	switch next {
	case tokens.Identifier, tokens.Symbol:
		p.pseudo = p.scanner.PseudoToken()
	default:
		p.pseudo = tokens.InvalidPseudoToken
	}
	return p.current
}

func (p *parser) expectIdent() ast.Name {
	p.builder.PushContext()
	defer p.builder.PopContext()
	if p.current == tokens.Identifier {
		result := p.builder.Name(p.scanner.Value().(string))
		p.next()
		return result
	}
	result := p.builder.Name("<error>")
	p.expect(tokens.Identifier)
	return result
}

func (p *parser) preserve() *parser {
	scanner := p.scanner.Clone()
	builder := p.builder.Clone(scanner)
	return &parser{scanner: scanner, builder: builder, current: p.current, pseudo: p.pseudo, operator: p.operator,
		separatorState: p.separatorState, errors: p.errors}
}

func (p *parser) restore(parser *parser) {
	p.scanner = parser.scanner
	p.builder = parser.builder
	p.current = parser.current
	p.pseudo = parser.pseudo
	p.operator = parser.operator
	p.separatorState = parser.separatorState
	p.errors = parser.errors
}

func (p *parser) firstOf(options ...func() ast.Element) ast.Element {
	preserved := p.preserve()
	firstErrorIndex := len(p.errors)
	var longestErrorOption *parser
	longestErrorEnd := location.Pos(0)
	var errorResult ast.Element
	for _, option := range options {
		result := option()
		if len(p.errors) > firstErrorIndex {
			e := p.errors[firstErrorIndex].End()
			if e > longestErrorEnd {
				longestErrorOption = p.preserve()
				errorResult = result
				longestErrorEnd = e
			}
			p.restore(preserved)
		} else {
			return result
		}
	}
	assert.Assert(longestErrorOption != nil, "An error option was expected")
	p.restore(longestErrorOption)
	return errorResult
}

func (p *parser) firstOfArray(options ...func() []ast.Element) []ast.Element {
	preserved := p.preserve()
	firstErrorIndex := len(p.errors)
	var longestErrorOption *parser
	longestErrorEnd := location.Pos(0)
	var errorResult []ast.Element
	for _, option := range options {
		result := option()
		if len(p.errors) > firstErrorIndex {
			e := p.errors[firstErrorIndex].End()
			if e > longestErrorEnd {
				longestErrorOption = p.preserve()
				errorResult = result
				longestErrorEnd = e
			}
			p.restore(preserved)
		} else {
			return result
		}
	}
	assert.Assert(longestErrorOption != nil, "An error option was expected")
	p.restore(longestErrorOption)
	return errorResult
}

func (p *parser) separator() bool {
	if p.current == tokens.Comma {
		p.next()
		return true
	} else if p.scanner.NewLineLocation().IsValid() {
		switch p.separatorState {
		case wasInfixState, separatorImplied:
			return false
		}
		if p.pseudo != tokens.Escaped {
			op := p.findOperator(ast.Infix, false)
			if op != nil {
				return false
			}
		}
		p.separatorState = separatorImplied
		return true
	}
	return false
}

var primitiveTokens = []tokens.Token{
	tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
	tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
	tokens.LBrace, tokens.LParen, tokens.Symbol, tokens.Let, tokens.LBrack, tokens.LBrackBang, tokens.LBraceBang,
}

func (p *parser) sequence() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()

	var left ast.Element
	switch p.current {
	case tokens.Identifier:
		switch p.pseudo {
		case tokens.Break:
			left = p.breakStatement()
		case tokens.Continue:
			left = p.continueStatement()
		case tokens.Loop:
			left = p.loopStatement()
		default:
			left = p.expression()
		}
	case tokens.Symbol:
		switch p.pseudo {
		case tokens.Spread:
			left = p.spreadExpression()
		case tokens.LessThan:
			left = p.typeLiteral()
		default:
			left = p.expression()
		}
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.LBrace,
		tokens.LParen, tokens.Let, tokens.LBrack, tokens.LBrackBang, tokens.LBraceBang:
		left = p.expression()
	case tokens.Var, tokens.Val:
		left = p.varDeclaration()
	case tokens.Return:
		left = p.returnStatement()
	default:
		left = p.expects(primitiveTokens...)
	}
	if p.separator() {
		switch p.current {
		case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
			tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LBrace,
			tokens.LParen, tokens.Symbol, tokens.Let, tokens.LBrack, tokens.LBrackBang, tokens.LBraceBang, tokens.Var, tokens.Val,
			tokens.Return:
			right := p.sequence()
			return p.builder.Sequence(left, right)
		}
	}
	return left
}

func (p *parser) expression() ast.Element {
	switch p.current {
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LBrace,
		tokens.LParen, tokens.Symbol, tokens.Let, tokens.LBrack, tokens.LBrackBang, tokens.LBraceBang:
		return p.operatorExpression(p.embeddingContext.lowestLevel)
	default:
		return p.expects(primitiveTokens...)
	}
}

type selectedOperator struct {
	name      ast.Name
	level     precedenceLevel
	assoc     ast.OperatorAssociativity
	placement ast.OperatorPlacement
}

func (op *selectedOperator) String() string {
	return fmt.Sprintf("op(%s, %d, %s, %s)", op.name, op.level.Level(), op.assoc, op.placement)
}

func selectOp(name ast.Name, op operator, placement ast.OperatorPlacement) *selectedOperator {
	level := op.Levels()[placement]
	if level != nil {
		return &selectedOperator{name: name, level: level, assoc: op.Associativities()[placement], placement: placement}
	}
	return nil
}

func (op *selectedOperator) isHigher(level precedenceLevel) bool {
	return (op.level == level && op.assoc == ast.Right) || op.level.IsHigherThan(level)
}

var noOperatorSentinal = &selectedOperator{}

func (p *parser) pushExcludeOperator(name string) {
	element, _ := p.vocabulary.Get(name)
	operator, _ := element.(operator)
	p.excludedOperators = append(p.excludedOperators, operator)
}

func (p *parser) popExcludedOperators() {
	p.excludedOperators = p.excludedOperators[0 : len(p.excludedOperators)-1]
}

func (p *parser) excludedOperator(operator operator) bool {
	for _, excludedOp := range p.excludedOperators {
		if excludedOp != nil && excludedOp == operator {
			return true
		}
	}
	return false
}

func (p *parser) findOperator(placement ast.OperatorPlacement, includeTypeMember bool) *selectedOperator {
	if placement == ast.Infix && p.operator != nil {
		if p.operator == noOperatorSentinal {
			return nil
		}
		return p.operator
	}
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.Identifier:
		if p.pseudo == tokens.Escaped {
			break
		}
		fallthrough
	case tokens.Symbol:
		text := p.scanner.Value().(string)
		element, ok := p.vocabulary.Get(text)
		if !ok {
			if includeTypeMember && placement == ast.Infix && p.current == tokens.Identifier {
				element, ok = p.vocabulary.Get(infixTypeMember)
				if ok {
				}
			}
			if !ok {
				if placement == ast.Infix {
					p.operator = noOperatorSentinal
				}
				return nil
			}
		}
		op, ok := element.(operator)
		if !ok {
			if placement == ast.Infix {
				p.operator = noOperatorSentinal
			}
			return nil
		}
		if p.excludedOperator(op) {
			p.operator = noOperatorSentinal
			return nil
		}
		if placement == ast.Postfix {
			text = "postfix " + text
		}
		name := p.builder.Name(text)
		p.operator = selectOp(name, op, placement)
		return p.operator
	}
	p.operator = noOperatorSentinal
	return nil
}

func (p *parser) unaryOp(target ast.Element, o *selectedOperator) ast.Element {
	return p.builder.Call(p.builder.Selection(target, o.name), nil)
}

func (p *parser) binaryOp(target ast.Element, o *selectedOperator, right ast.Element) ast.Element {
	return p.builder.Call(p.builder.Selection(target, o.name), []ast.Element{right})
}

func (p *parser) operatorExpression(level precedenceLevel) ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()

	var left ast.Element
	op := p.findOperator(ast.Prefix, false)
	if op != nil && op.isHigher(level) {
		p.next()
		left = p.unaryOp(p.operatorExpression(op.level), op)
	} else {
		left = p.simpleExpression()
	}
	op = p.findOperator(ast.Postfix, false)
	for op != nil && op.isHigher(level) {
		p.next()
		left = p.unaryOp(left, op)
		op = p.findOperator(ast.Postfix, false)
	}
	op = p.findOperator(ast.Infix, !p.scanner.NewLineLocation().IsValid())
	for op != nil && op.isHigher(level) {
		p.next()
		p.separatorState = wasInfixState
		right := p.operatorExpression(op.level)
		left = p.binaryOp(left, op, right)
		op = p.findOperator(ast.Infix, true)
	}
	return left
}

func (p *parser) simpleExpression() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LBrace,
		tokens.LParen, tokens.Let, tokens.LBrack, tokens.LBrackBang, tokens.LBraceBang:
		left := p.primitive()
		for {
			switch p.current {
			case tokens.Dot:
				left = p.selector(left)
				continue
			case tokens.LParen:
				left = p.call(left)
				continue
			case tokens.LBrack:
				left = p.index(left)
				continue
			}
			break
		}
		return left
	case tokens.Symbol:
		text := p.scanner.Value()
		p.next()
		return p.report("Symbol '%s' is not defined as an operator in the current vocabulary", text)
	default:
		return p.expects(primitiveTokens...)
	}
}

func (p *parser) selector(left ast.Element) ast.Element {
	p.expect(tokens.Dot)
	name := p.expectIdent()
	return p.builder.Selection(left, name)
}

func (p *parser) call(left ast.Element) ast.Element {
	p.expect(tokens.LParen)
	arguments := p.arguments()
	p.expect(tokens.RParen)
	return p.builder.Call(left, arguments)
}

func (p *parser) index(left ast.Element) ast.Element {
	p.expect(tokens.LBrack)
	arguments := p.arguments()
	p.expect(tokens.RBrack)
	if p.pseudo == tokens.Equal {
		p.next()
		arguments = append(arguments, p.expression())
		name := p.builder.Name("set")
		selection := p.builder.Selection(left, name)
		return p.builder.Call(selection, arguments)
	}
	name := p.builder.Name("get")
	selection := p.builder.Selection(left, name)
	return p.builder.Call(selection, arguments)
}

func (p *parser) arguments() []ast.Element {
	var arguments []ast.Element
	switch p.current {
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
		tokens.LBrace, tokens.LParen, tokens.Symbol, tokens.Colon:
		for {
			switch p.current {
			case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
				tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
				tokens.LBrace, tokens.LParen, tokens.Symbol, tokens.Colon:
				argument := p.argument()
				arguments = append(arguments, argument)
				if p.separator() {
					switch p.current {
					case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
						tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
						tokens.LBrace, tokens.LParen, tokens.Symbol, tokens.Colon:
						continue
					}
				}
			}
			break
		}
	}
	return arguments
}

func (p *parser) argument() ast.Element {
	return p.firstOf(func() ast.Element {
		return p.namedArgument()
	}, func() ast.Element {
		result := p.expression()
		if p.current == tokens.Colon {
			p.expect(tokens.Comma)
		}
		return result
	})
}

func (p *parser) namedArgument() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	var name ast.Name
	if p.current == tokens.Colon {
		p.next()
		if p.current != tokens.Identifier {
			name = p.expectIdent()
		} else {
			name = p.builder.Name(p.scanner.Value().(string))
		}
	} else {
		name = p.expectIdent()
		p.expect(tokens.Colon)
	}
	value := p.expression()
	return p.builder.NamedArgument(name, value)
}

func (p *parser) whenExpression() ast.When {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.When)
	var target ast.Element
	if p.current == tokens.LParen {
		p.next()
		target = p.expression()
		p.expect(tokens.RParen)
	}
	p.expect(tokens.LBrace)
	clauses := p.whenClauses()
	p.expect(tokens.RBrace)
	p.builder.PushContext()
	return p.builder.When(target, clauses)
}

func (p *parser) whenClauses() []ast.Element {
	var result []ast.Element
	for {
		switch p.current {
		case tokens.Identifier:
			if p.pseudo == tokens.Else {
				result = append(result, p.whenElseClause())
				if p.separator() {
					continue
				} else {
					break
				}
			}
			fallthrough
		case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
			tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False,
			tokens.LBrace, tokens.LParen, tokens.Symbol:
			result = append(result, p.whenValueClause())
			if p.separator() {
				continue
			}
		}
		if p.separator() {
			switch p.current {
			case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
				tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
				tokens.LBrace, tokens.LParen, tokens.Symbol:
				continue
			}
		}
		break
	}
	p.separator()
	return result
}

func (p *parser) whenElseClause() ast.WhenElseClause {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.Else)
	p.expectPseudo(tokens.Arrow)
	p.expect(tokens.LBrace)
	body := p.sequence()
	p.expect(tokens.RBrace)
	return p.builder.WhenElseClause(body)
}

func (p *parser) whenValueClause() ast.WhenValueClause {
	p.builder.PushContext()
	defer p.builder.PopContext()
	value := p.expression()
	p.expectPseudo(tokens.Arrow)
	p.expect(tokens.LBrace)
	body := p.sequence()
	p.expect(tokens.RBrace)
	return p.builder.WhenValueClause(value, body)
}

func (p *parser) lambda() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.LBrace)
	typeParameters := p.typeParametersClause()
	parameters := p.lambdaParameters()
	var expression ast.Element
	if p.current != tokens.RBrace {
		expression = p.sequence()
	}
	p.expect(tokens.RBrace)
	return p.builder.Lambda(typeParameters, parameters, expression)
}

func (p *parser) intrinsicLambda() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.LBraceBang)
	typeParameters := p.typeParametersClause()
	parameters := p.lambdaParameters()
	var sequence ast.Element
	if p.current != tokens.BangRBrace {
		sequence = p.sequence()
	}
	p.expect(tokens.BangRBrace)
	var resultType ast.Element
	if p.current == tokens.Colon {
		p.next()
		resultType = p.typeReference()
	}
	return p.builder.IntrinsicLambda(typeParameters, parameters, sequence, resultType)
}

func (p *parser) typeParametersClause() ast.TypeParameters {
	result, _ := p.firstOf(func() ast.Element {
		p.builder.PushContext()
		defer p.builder.PopContext()
		typeParameters := p.typeParameters()
		whereClauses := p.whereClauses()
		p.expectPseudo(tokens.Bar)
		return p.builder.TypeParameters(typeParameters, whereClauses)
	}, func() ast.Element {
		return nil
	}).(ast.TypeParameters)
	return result
}

func (p *parser) typeParameters() []ast.TypeParameter {
	result := p.firstOfArray(func() []ast.Element {
		var result []ast.Element
		for {
			switch p.current {
			case tokens.Identifier:
				typeParameter := p.typeParameter()
				result = append(result, typeParameter)
				if p.separator() {
					if p.pseudo != tokens.Bar && p.pseudo != tokens.Where {
						continue
					}
				}
			case tokens.Symbol:
				if p.pseudo == tokens.Bar {
					break
				}
			}
			break
		}
		switch p.pseudo {
		case tokens.Bar, tokens.Where:
			return result
		}
		p.expectPseudo(tokens.Bar)
		return result
	}, func() []ast.Element {
		return nil
	})
	var params []ast.TypeParameter
	for _, param := range result {
		params = append(params, param.(ast.TypeParameter))
	}
	return params
}

func (p *parser) whereClauses() []ast.Where {
	var result []ast.Where
	for {
		if p.pseudo == tokens.Where {
			whereClause := p.whereClause()
			result = append(result, whereClause)
			continue
		}
		break
	}
	return result
}

func (p *parser) whereClause() ast.Where {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.Where)
	left := p.typeReference()
	p.expectPseudo(tokens.Equal)
	right := p.typeReference()
	return p.builder.Where(left, right)
}

func (p *parser) typeParameter() ast.TypeParameter {
	p.builder.PushContext()
	defer p.builder.PopContext()
	name := p.expectIdent()
	var typeReference ast.Element
	if p.current == tokens.Colon {
		p.expect(tokens.Colon)
		typeReference = p.typeReference()
	}
	return p.builder.TypeParameter(name, typeReference)
}

func (p *parser) lambdaParameters() []ast.Parameter {
	result := p.firstOfArray(func() []ast.Element {
		result := p.parameters()
		p.expectPseudo(tokens.Arrow)
		return result
	}, func() []ast.Element {
		return nil
	})
	var params []ast.Parameter
	for _, param := range result {
		params = append(params, param.(ast.Parameter))
	}
	return params
}

func (p *parser) parameters() []ast.Element {
	var result []ast.Element
	for {
		if p.current == tokens.Identifier {
			parameter := p.parameter()
			result = append(result, parameter)
			if p.separator() {
				if p.current == tokens.Identifier {
					continue
				}
			}
		}
		break
	}
	return result
}

func (p *parser) parameter() ast.Parameter {
	p.builder.PushContext()
	defer p.builder.PopContext()
	name := p.expectIdent()
	var typeReference ast.Element
	if p.current == tokens.Colon {
		p.next()
		typeReference = p.typeReference()
	}
	var defaultExpression ast.Element
	if p.pseudo == tokens.Equal {
		p.next()
		defaultExpression = p.expression()
	}
	return p.builder.Parameter(name, typeReference, defaultExpression)
}

func (p *parser) primitive() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.LiteralRune:
		result := p.builder.LiteralRune(p.scanner.Value().(rune))
		p.next()
		return result
	case tokens.LiteralByte:
		result := p.builder.LiteralByte(p.scanner.Value().(byte))
		p.next()
		return result
	case tokens.LiteralInt:
		result := p.builder.LiteralInt(p.scanner.Value().(int))
		p.next()
		return result
	case tokens.LiteralUInt:
		result := p.builder.LiteralUInt(p.scanner.Value().(uint))
		p.next()
		return result
	case tokens.LiteralLong:
		result := p.builder.LiteralLong(p.scanner.Value().(int64))
		p.next()
		return result
	case tokens.LiteralULong:
		result := p.builder.LiteralULong(p.scanner.Value().(uint64))
		p.next()
		return result
	case tokens.LiteralDouble:
		result := p.builder.LiteralDouble(p.scanner.Value().(float64))
		p.next()
		return result
	case tokens.LiteralFloat:
		result := p.builder.LiteralFloat(p.scanner.Value().(float32))
		p.next()
		return result
	case tokens.LiteralString:
		result := p.builder.LiteralString(p.scanner.Value().(string))
		p.next()
		return result
	case tokens.True:
		result := p.builder.LiteralBoolean(true)
		p.next()
		return result
	case tokens.False:
		result := p.builder.LiteralBoolean(false)
		p.next()
		return result
	case tokens.Identifier:
		switch p.pseudo {
		case tokens.When:
			return p.whenExpression()
		}
		result := p.builder.Name(p.scanner.Value().(string))
		p.next()
		return result
	case tokens.LBrace:
		return p.lambda()
	case tokens.LBraceBang:
		return p.intrinsicLambda()
	case tokens.LBrack:
		return p.firstOf(func() ast.Element {
			return p.readOnlyObjectInitializer()
		}, func() ast.Element {
			return p.readOnlyArrayInitializer()
		})
	case tokens.LBrackBang:
		return p.firstOf(func() ast.Element {
			return p.mutableObjectInitializer()
		}, func() ast.Element {
			return p.mutableArrayInitializer()
		})
	case tokens.LParen:
		excludedOperators := p.excludedOperators
		p.excludedOperators = nil
		p.expect(tokens.LParen)
		expr := p.expression()
		p.expect(tokens.RParen)
		p.excludedOperators = excludedOperators
		return expr
	case tokens.Let:
		return p.definition()
	default:
		return p.expects(primitiveTokens...)
	}
}

func (p *parser) spreadExpression() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.Spread)
	var target ast.Element
	var vocabulary vocabulary
	if p.current == tokens.VocabularyStart {
		vocabularyLiteral := p.vocabularyLiteral()
		v, errors := buildVocabulary(p.scope, vocabularyLiteral)
		vocabulary = v
		if errors != nil {
			for _, err := range errors {
				p.reportElement(err.element, err.message)
			}
		}
		target = vocabularyLiteral
	} else {
		preserved := p.preserve()
		target := p.spreadReference()
		if len(p.errors) > len(preserved.errors) {
			p.restore(preserved)
			target := p.expression()
			return p.builder.Spread(target)
		}

		// Find an apply vocabulary
		v, element := lookupVocabulary(p.scope, target)
		vocabulary = v
		if element != nil {
			return p.reportElement(element, "Expected a vocabulary reference")
		}
	}
	p.embeddingContext.embedVocabulary(vocabulary.(*vocabularyImpl), target)
	for _, err := range p.embeddingContext.errors {
		p.reportElement(target, err.message)
	}
	p.embeddingContext.errors = nil
	return p.builder.Spread(target)
}

func (p *parser) spreadReference() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	var result ast.Element = p.expectIdent()
	for p.current == tokens.Scope {
		p.next()
		name := p.expectIdent()
		result = p.builder.Selection(result, name)
	}
	return result
}

func (p *parser) typeReference() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	result := p.simpleTypeReference()
	operator := p.typeOperator()
	if operator != nil {
		target := p.builder.Selection(result, operator)
		result = p.builder.Call(target, []ast.Element{p.typeReference()})
	}
	return result
}

func (p *parser) typeOperator() ast.Name {
	p.builder.PushContext()
	defer p.builder.PopContext()
	if p.current == tokens.Symbol {
		if p.pseudo == tokens.And {
			name := p.builder.Name(p.scanner.Value().(string))
			p.next()
			return name
		}
	}
	return nil
}

func (p *parser) simpleTypeReference() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	result := p.typeReferencePrimitive()
	for {
		switch p.current {
		case tokens.Dot:
			p.next()
			name := p.expectIdent()
			result = p.builder.Selection(result, name)
			continue
		case tokens.LBrack:
			p.next()
			p.expect(tokens.RBrack)
			result = p.builder.SequenceType(result)
			continue
		case tokens.Symbol:
			if p.pseudo == tokens.Question {
				p.next()
				result = p.builder.OptionalType(result)
				continue
			}
		}
		break
	}
	return result
}

func (p *parser) typeReferencePrimitive() ast.Element {
	switch p.current {
	case tokens.LParen:
		p.next()
		result := p.typeReference()
		p.expect(tokens.RParen)
		return result
	case tokens.Symbol:
		if p.pseudo == tokens.LessThan {
			return p.typeLiteral()
		}
		fallthrough
	default:
		return p.expectIdent()
	}
}

func (p *parser) typeLiteral() ast.TypeLiteral {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.LessThan)
	p.pushExcludeOperator(tokens.GreaterThan.String())
	defer p.popExcludedOperators()
	var members []ast.Element
	for {
		switch p.current {
		case tokens.LBrace:
			members = append(members, p.callableTypeMember())
		case tokens.Symbol:
			if p.pseudo == tokens.GreaterThan {
				break
			} else if p.pseudo == tokens.Spread {
				members = append(members, p.spreadTypeMember())
				continue
			}
			fallthrough
		case tokens.Identifier:
			members = append(members, p.typeLiteralMember())
		}
		if p.separator() {
			continue
		}
		break
	}
	p.separator()
	p.expectPseudo(tokens.GreaterThan)
	return p.builder.TypeLiteral(members)
}

func (p *parser) spreadTypeMember() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.Spread)
	reference := p.typeReference()
	return p.builder.SpreadTypeMember(reference)
}

func (p *parser) typeLiteralMember() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	var name ast.Name
	switch p.current {
	case tokens.Identifier, tokens.Symbol:
		name = p.builder.Name(p.scanner.Value().(string))
		p.next()
	default:
		name = p.expectIdent()
	}
	switch p.current {
	case tokens.Colon:
		p.next()
		typ := p.typeReference()
		return p.builder.TypeLiteralMember(name, typ)
	case tokens.Symbol:
		p.expectPseudo(tokens.Equal)
		var value ast.Element
		if p.pseudo == tokens.LessThan {
			value = p.typeLiteral()
		} else {
			value = p.expression()
		}
		return p.builder.TypeLiteralConstant(name, value)
	default:
		return p.expects(tokens.Colon)
	}
}

func (p *parser) callableTypeMember() ast.CallableTypeMember {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.LBrace)
	parameters := p.parameters()
	p.expectPseudo(tokens.Arrow)
	resultType := p.typeReference()
	p.expect(tokens.RBrace)
	return p.builder.CallableTypeMember(parameters, resultType)
}

func (p *parser) definition() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.Let:
		p.next()
		name := p.definitionName()
		p.expectPseudo(tokens.Equal)
		value := p.letValue()
		return p.builder.LetDefinition(name, value)
	default:
		return p.expects(tokens.Let, tokens.Var)
	}
}

func (p *parser) definitionName() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	var result ast.Element = p.expectIdent()
	for p.current == tokens.Dot {
		p.next()
		name := p.expectIdent()
		result = p.builder.Selection(result, name)
	}
	return result
}

func (p *parser) letValue() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.LBrace:
		return p.lambda()
	case tokens.VocabularyStart:
		return p.vocabularyLiteral()
	case tokens.LiteralString, tokens.LiteralRune, tokens.LiteralInt, tokens.LiteralByte, tokens.LiteralUInt,
		tokens.LiteralLong, tokens.LiteralULong, tokens.LiteralFloat, tokens.LiteralDouble:
		return p.primitive()
	case tokens.Symbol:
		if p.pseudo == tokens.LessThan {
			return p.typeLiteral()
		}
		fallthrough
	default:
		return p.expectItems(tokens.VocabularyStart, tokens.LessThan, tokens.LBrace)
	}
}

func (p *parser) varDeclaration() ast.VarDefinition {
	p.builder.PushContext()
	defer p.builder.PopContext()
	var mutable = false
	if p.current == tokens.Var {
		mutable = true
		p.next()
	} else {
		p.expect(tokens.Val)
	}
	name := p.expectIdent()
	var typ ast.Element
	if p.current == tokens.Colon {
		p.next()
		typ = p.typeReference()
	}
	var value ast.Element
	if p.pseudo == tokens.Equal {
		p.next()
		value = p.expression()
	}
	return p.builder.VarDefinition(name, typ, value, mutable)
}

func (p *parser) readOnlyObjectInitializer() ast.ObjectInitializer {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.LBrack)
	typ := p.initializerType()
	members := p.memberInitializers()
	p.expect(tokens.RBrack)
	return p.builder.ObjectInitializer(false, typ, members)
}

func (p *parser) initializerType() ast.Element {
	if p.pseudo == tokens.LessThan {
		p.next()
		typ := p.typeReference()
		p.expectPseudo(tokens.GreaterThan)
		return typ
	}
	return nil
}

func (p *parser) mutableObjectInitializer() ast.ObjectInitializer {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.LBrackBang)
	typ := p.initializerType()
	members := p.memberInitializers()
	p.expect(tokens.BangRBrack)
	return p.builder.ObjectInitializer(true, typ, members)
}

func (p *parser) memberInitializers() []ast.Element {
	var result []ast.Element
	for {
		switch p.current {
		case tokens.Identifier, tokens.Colon:
			result = append(result, p.memberInitializer())
		case tokens.Symbol:
			if p.pseudo == tokens.Spread {
				result = append(result, p.memberInitializer())
			}
		}
		if p.separator() {
			continue
		}
		break
	}
	return result
}

func (p *parser) memberInitializer() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.Colon:
		p.next()
		var name ast.Name
		if p.current == tokens.Identifier {
			name = p.builder.Name(p.scanner.Value().(string))
		} else {
			name = p.expectIdent()
		}
		// Intentionally do not call p.next() so the identifier is considered part of the expression
		value := p.expression()
		return p.builder.NamedMemberInitializer(name, nil, value)
	case tokens.Identifier:
		name := p.expectIdent()
		p.expect(tokens.Colon)
		value := p.expression()
		return p.builder.NamedMemberInitializer(name, nil, value)
	case tokens.Symbol:
		if p.pseudo == tokens.Spread {
			p.next()
			spreadValue := p.expression()
			return p.builder.SpreadMemberInitializer(spreadValue)
		}
	}
	return p.expectIdent()
}

func (p *parser) readOnlyArrayInitializer() ast.ArrayInitializer {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.LBrack)
	typ := p.initializerType()
	elements := p.arrayElements()
	p.expect(tokens.RBrack)
	return p.builder.ArrayInitializer(false, typ, elements)
}

func (p *parser) mutableArrayInitializer() ast.ArrayInitializer {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.LBrackBang)
	typ := p.initializerType()
	elements := p.arrayElements()
	p.expect(tokens.BangRBrack)
	return p.builder.ArrayInitializer(true, typ, elements)
}

func (p *parser) arrayElements() []ast.Element {
	var elements []ast.Element
	for {
		switch p.current {
		case tokens.Symbol:
			if p.pseudo == tokens.Spread {
				elements = append(elements, p.memberInitializer())
				if p.separator() {
					continue
				}
				break
			}
			fallthrough
		case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
			tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LBrace,
			tokens.LParen, tokens.Let, tokens.LBrack, tokens.LBrackBang, tokens.LBraceBang:
			elements = append(elements, p.expression())
		}
		if p.separator() {
			continue
		}
		break
	}
	p.separator()
	return elements
}

func (p *parser) vocabularyLiteral() ast.VocabularyLiteral {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.VocabularyStart)
	members := p.vocabularyMembers()
	p.expect(tokens.VocabularyEnd)
	return p.builder.VocabularyLiteral(members)
}

func (p *parser) vocabularyMembers() []ast.Element {
	var result []ast.Element
	for {
		switch p.current {
		case tokens.Identifier:
			switch p.pseudo {
			case tokens.Infix, tokens.Prefix, tokens.Postfix:
				operator := p.vocabularyOperatorDeclaration()
				result = append(result, operator)
				continue
			default:
				p.expectsPseudo(tokens.Infix, tokens.Prefix, tokens.Postfix)
				continue
			}
		case tokens.Symbol:
			if p.pseudo == tokens.Spread {
				result = append(result, p.vocabularyEmbedding())
			} else {
				break
			}
		default:
			if p.separator() {
				continue
			}
		}
		break
	}
	return result
}

func (p *parser) vocabularyOperatorDeclaration() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	placement := ast.Infix
	switch p.pseudo {
	case tokens.Infix:
		placement = ast.Infix
		p.next()
	case tokens.Prefix:
		placement = ast.Prefix
		p.next()
	case tokens.Postfix:
		placement = ast.Postfix
		p.next()
	default:
		return p.expectsPseudo(tokens.Infix, tokens.Prefix, tokens.Postfix)
	}
	p.expectPseudo(tokens.Operator)
	names := p.names()
	qualifier := p.vocabularyPrecedenceQualifier()
	associativity := ast.Left
	switch p.pseudo {
	case tokens.Left:
		associativity = ast.Left
		p.next()
	case tokens.Right:
		associativity = ast.Right
		p.next()
	default:
		return p.expectsPseudo(tokens.Left, tokens.Right)
	}
	return p.builder.VocabularyOperatorDeclaration(names, placement, qualifier, associativity)
}

func (p *parser) names() []ast.Name {
	switch p.current {
	case tokens.Identifier:
		var name ast.Name
		if p.pseudo == tokens.Identifiers {
			name = p.builder.Name(infixTypeMember)
			p.next()
		} else {
			name = p.expectIdent()
		}
		return []ast.Name{name}
	case tokens.LParen:
		p.next()
		var result []ast.Name
		for {
			switch p.current {
			case tokens.Identifier:
				name := p.expectIdent()
				result = append(result, name)
				fallthrough
			default:
				if p.separator() {
					continue
				}
			}
			break
		}
		p.expect(tokens.RParen)
		return result
	default:
		p.expects(tokens.Identifier, tokens.LParen)
		return nil
	}
}

func (p *parser) vocabularyPrecedenceQualifier() ast.VocabularyOperatorPrecedence {
	p.builder.PushContext()
	defer p.builder.PopContext()
	relation := ast.After
	switch p.pseudo {
	case tokens.After:
		relation = ast.After
		p.next()
	case tokens.Before:
		relation = ast.Before
		p.next()
	default:
		return nil
	}
	placement := ast.UnspecifiedPlacement
	switch p.pseudo {
	case tokens.Infix:
		placement = ast.Infix
		p.next()
	case tokens.Prefix:
		placement = ast.Prefix
		p.next()
	case tokens.Postfix:
		placement = ast.Postfix
		p.next()
	}
	name := p.expectIdent()
	return p.builder.VocabularyOperatorPrecedence(name, placement, relation)
}

func (p *parser) vocabularyEmbedding() ast.VocabularyEmbedding {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.Spread)
	name := p.vocabularyNameReference()
	return p.builder.VocabularyEmbedding(name)
}

func (p *parser) vocabularyNameReference() []ast.Name {
	var result []ast.Name
	for {
		name := p.expectIdent()
		result = append(result, name)
		if p.current == tokens.Scope {
			p.next()
			continue
		}
		break
	}
	return result
}

func lookupVocabulary(scope vocabularyScope, element ast.Element) (vocabulary, ast.Element) {
	result, elem := lookup(scope, element)
	if elem != nil {
		return nil, elem
	}
	vocab, ok := result.(vocabulary)
	if !ok {
		return nil, element
	}
	return vocab, nil
}

func lookup(scope vocabularyScope, element ast.Element) (any, ast.Element) {
	switch e := element.(type) {
	case ast.Name:
		result, ok := scope.Get(e.Text())
		if !ok {
			return nil, e
		}
		return result, nil
	case ast.Selection:
		sc, elem := lookup(scope, e.Target())
		if elem != nil {
			return nil, elem
		}
		newScope, ok := sc.(vocabularyScope)
		if !ok {
			return nil, e
		}
		result, elem := lookup(newScope, e.Member())
		return result, elem
	default:
		assert.Fail("selectorsToNames called with invalid node %#v", element)
	}
	return nil, element
}

func (p *parser) loopStatement() ast.Loop {
	p.builder.PushContext()
	defer p.builder.PopContext()
	var label ast.Name
	p.expectPseudo(tokens.Loop)
	if p.current == tokens.Identifier {
		label = p.expectIdent()
	}
	p.expect(tokens.LBrace)
	body := p.sequence()
	p.expect(tokens.RBrace)
	return p.builder.Loop(label, body)
}

func (p *parser) breakStatement() ast.Break {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.Break)
	var label ast.Name
	if p.current == tokens.Identifier {
		label = p.expectIdent()
	}
	return p.builder.Break(label)
}

func (p *parser) continueStatement() ast.Continue {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.Continue)
	var label ast.Name
	if p.current == tokens.Identifier {
		label = p.expectIdent()
	}
	return p.builder.Continue(label)
}

func (p *parser) returnStatement() ast.Return {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expect(tokens.Return)
	var value ast.Element
	switch p.current {
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False,
		tokens.LBrace, tokens.LParen, tokens.Symbol, tokens.LBrack, tokens.LBrackBang, tokens.LBraceBang, tokens.Identifier:
		value = p.expression()
	}
	return p.builder.Return(value)
}
