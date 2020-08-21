package parser

import (
	"fmt"

	"dyego0/assert"
	"dyego0/ast"
	"dyego0/scanner"
	"dyego0/tokens"
)

// Parser parses text and returns an ast element
type Parser interface {
	Errors() []ast.Error
	Parse() ast.Element
}

type parser struct {
	scanner          *scanner.Scanner
	builder          ast.Builder
	current          tokens.Token
	pseudo           tokens.PseudoToken
	scope            vocabularyScope
	vocabulary       vocabulary
	embeddingContext *vocabularyEmbeddingContext
	errors           []ast.Error
}

// NewParser creates a new parser
func NewParser(scanner *scanner.Scanner, scope vocabularyScope) Parser {
	builder := ast.NewBuilder(scanner)
	context := newVocabularyEmbeddingContext()
	p := &parser{
		scanner:          scanner,
		builder:          builder,
		pseudo:           tokens.InvalidPseudoToken,
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

func (p *parser) Errors() []ast.Error {
	return p.errors
}

func (p *parser) report(msg string, args ...interface{}) ast.Error {
	err := p.builder.Error(fmt.Sprintf(msg, args...))
	p.errors = append(p.errors, err)
	return err
}

func (p *parser) reportElement(element ast.Element, msg string, args ...interface{}) ast.Element {
	err := p.builder.DirectError(element.Start(), element.End(), fmt.Sprintf(msg, args...))
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

func (p *parser) next() tokens.Token {
	var next = p.scanner.Next()
	p.current = next
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
	return &parser{scanner: scanner, builder: p.builder.Clone(scanner), current: p.current, pseudo: p.pseudo, errors: p.errors}
}

func (p *parser) restore(parser *parser) {
	p.scanner = parser.scanner
	p.builder = parser.builder
	p.current = parser.current
	p.pseudo = parser.pseudo
	p.errors = parser.errors
}

var primitiveTokens = []tokens.Token{
	tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
	tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
	tokens.LBrace, tokens.LParen, tokens.Symbol, tokens.Let,
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
		if p.pseudo == tokens.Spread {
			left = p.spreadExpression()
		} else {
			left = p.expression()
		}
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.LBrace,
		tokens.LParen, tokens.Let:
		left = p.expression()
	case tokens.Return:
		left = p.returnStatement()
	default:
		left = p.expects(primitiveTokens...)
	}
	if p.current == tokens.Comma {
		p.next()
		switch p.current {
		case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
			tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LBrace,
			tokens.LParen, tokens.Symbol, tokens.Let:
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
		tokens.LParen, tokens.Symbol, tokens.Let:
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

func selectOp(name ast.Name, op operator, placement ast.OperatorPlacement) *selectedOperator {
	level := op.Levels()[placement]
	if level != nil {
		return &selectedOperator{name: name, level: level, assoc: op.Associativities()[placement], placement: placement}
	}
	return nil
}

func (o *selectedOperator) isHigher(level precedenceLevel) bool {
	return (o.level == level && o.assoc == ast.Right) || o.level.IsHigherThan(level)
}

func (p *parser) findOperator(placement ast.OperatorPlacement) *selectedOperator {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.Identifier, tokens.Symbol:
		text := p.scanner.Value().(string)
		element, ok := p.vocabulary.Get(text)
		if !ok {
			return nil
		}
		op, ok := element.(operator)
		if !ok {
			return nil
		}
		if placement == ast.Postfix {
			text = "postfix " + text
		}
		name := p.builder.Name(text)
		return selectOp(name, op, placement)
	}
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
	op := p.findOperator(ast.Prefix)
	if op != nil && op.isHigher(level) {
		p.next()
		left = p.unaryOp(p.operatorExpression(op.level), op)
	} else {
		left = p.simpleExpression()
	}
	op = p.findOperator(ast.Postfix)
	for op != nil && op.isHigher(level) {
		p.next()
		left = p.unaryOp(left, op)
		op = p.findOperator(ast.Postfix)
	}
	op = p.findOperator(ast.Infix)
	for op != nil && op.isHigher(level) {
		p.next()
		right := p.operatorExpression(op.level)
		left = p.binaryOp(left, op, right)
		op = p.findOperator(ast.Infix)
	}
	return left
}

func (p *parser) simpleExpression() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LBrace,
		tokens.LParen, tokens.Let:
		left := p.primitive()
		for {
			switch p.current {
			case tokens.Dot:
				left = p.selector(left)
				continue
			case tokens.LParen:
				left = p.call(left)
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

func (p *parser) arguments() []ast.Element {
	var arguments []ast.Element
	switch p.current {
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
		tokens.LBrace, tokens.LParen, tokens.Symbol:
		for {
			switch p.current {
			case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
				tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
				tokens.LBrace, tokens.LParen, tokens.Symbol:
				argument := p.argument()
				arguments = append(arguments, argument)
				if p.current == tokens.Comma {
					p.next()
					switch p.current {
					case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
						tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
						tokens.LBrace, tokens.LParen, tokens.Symbol:
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
	preserved := p.preserve()
	namedArgument := p.namedArgument()
	if len(p.errors) > len(preserved.errors) {
		// not a named arguemnt
		p.restore(preserved)
		return p.expression()
	}
	return namedArgument
}

func (p *parser) namedArgument() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	name := p.expectIdent()
	p.expectPseudo(tokens.Equal)
	value := p.expression()
	return p.builder.NamedArgument(name, value)
}

func (p *parser) whenExpression() ast.When {
	p.builder.PushContext()
	defer p.builder.PopContext()
	p.expectPseudo(tokens.When)
	p.expect(tokens.LParen)
	target := p.expression()
	p.expect(tokens.RParen)
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
				if p.current == tokens.Comma {
					p.next()
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
			if p.current == tokens.Comma {
				p.next()
				continue
			}
		}
		if p.current == tokens.Comma {
			switch p.current {
			case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
				tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
				tokens.LBrace, tokens.LParen, tokens.Symbol:
				continue
			}
		}
		break
	}
	if p.current == tokens.Comma {
		p.next()
	}
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
		expression = p.expression()
	}
	p.expect(tokens.RBrace)
	return p.builder.Lambda(typeParameters, parameters, expression)
}

func (p *parser) typeParametersClause() ast.TypeParameters {
	p.builder.PushContext()
	defer p.builder.PopContext()
	state := p.preserve()
	typeParameters := p.typeParameters()
	whereClauses := p.whereClauses()
	if p.pseudo != tokens.Bar {
		p.restore(state)
		return nil
	}
	p.expectPseudo(tokens.Bar)
	return p.builder.TypeParameters(typeParameters, whereClauses)
}

func (p *parser) typeParameters() []ast.TypeParameter {
	state := p.preserve()
	var result []ast.TypeParameter
	for {
		switch p.current {
		case tokens.Identifier:
			typeParameter := p.typeParameter()
			result = append(result, typeParameter)
			if p.current == tokens.Comma {
				p.expect(tokens.Comma)
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
	if p.pseudo != tokens.Bar && p.pseudo != tokens.Where {
		p.restore(state)
		return nil
	}
	return result
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
	state := p.preserve()
	result := p.parameters()
	if p.pseudo != tokens.Arrow {
		p.restore(state)
		return nil
	}
	p.expectPseudo(tokens.Arrow)
	return result
}

func (p *parser) parameters() []ast.Parameter {
	var result []ast.Parameter
	for {
		if p.current == tokens.Identifier {
			parameter := p.parameter()
			result = append(result, parameter)
			if p.current == tokens.Comma {
				p.next()
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
	case tokens.LParen:
		p.expect(tokens.LParen)
		expr := p.expression()
		p.expect(tokens.RParen)
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
	preserved := p.preserve()
	reference := p.spreadReference()
	if len(p.errors) > len(preserved.errors) {
		p.restore(preserved)
		target := p.expression()
		return p.builder.Spread(target)
	}

	// Find an apply vocabulary
	vocabulary, element := lookupVocabulary(p.scope, reference)
	if element != nil {
		return p.reportElement(element, "Expected a vocabulary reference")
	}
	p.embeddingContext.embedVocabulary(vocabulary.(*vocabularyImpl), reference)
	for _, err := range p.embeddingContext.errors {
		p.reportElement(reference, err.message)
	}
	p.embeddingContext.errors = nil
	return p.builder.Spread(reference)
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
	// TODO: Implement
	return p.expectIdent()
}

func (p *parser) definition() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.Let:
		p.next()
		name := p.expectIdent()
		p.expectPseudo(tokens.Equal)
		value := p.letValue()
		return p.builder.LetDefinition(name, value)
	default:
		return p.expects(tokens.Let, tokens.Var)
	}
}

func (p *parser) letValue() ast.Element {
	p.builder.PushContext()
	defer p.builder.PopContext()
	switch p.current {
	case tokens.VocabularyStart:
		return p.vocabularyLiteral()
	default:
		return p.expects(tokens.VocabularyStart)
	}
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
		case tokens.Comma:
			p.next()
			continue
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
		name := p.expectIdent()
		return []ast.Name{name}
	case tokens.LParen:
		p.next()
		var result []ast.Name
		for {
			switch p.current {
			case tokens.Identifier:
				name := p.expectIdent()
				result = append(result, name)
				continue
			case tokens.Comma:
				p.next()
				continue
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
		tokens.LBrace, tokens.LParen, tokens.Symbol:
		value = p.expression()
	}
	return p.builder.Return(value)
}
