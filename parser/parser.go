package parser

import (
	"fmt"

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
	scanner *scanner.Scanner
	builder ast.Builder
	current tokens.Token
	pseudo  tokens.PseudoToken
	errors  []ast.Error
}

// NewParser creates a new parser
func NewParser(scanner *scanner.Scanner) Parser {
	builder := ast.NewBuilder(scanner)
	p := &parser{scanner: scanner, builder: builder, pseudo: tokens.InvalidPseudoToken}
	builder.PushContext()
	p.next()
	return p
}

func (p *parser) Parse() ast.Element {
	expr := p.expression()
	p.expect(tokens.EOF)
	return expr
}

func (p *parser) Errors() []ast.Error {
	return p.errors
}

func (p *parser) report(msg string, args ...interface{}) {
	p.errors = append(p.errors, p.builder.Error(fmt.Sprintf(msg, args...)))
}

func (p *parser) expect(t tokens.Token) {
	if p.current == t {
		p.next()
	} else {
		p.builder.PushContext()
		defer p.builder.PopContext()
		p.report("Expected %v, recieved %v", t, p.current)
		p.next()
	}
}

func (p *parser) expectPseudo(t tokens.PseudoToken) {
	if p.pseudo == t {
		p.next()
	} else {
		p.builder.PushContext()
		defer p.builder.PopContext()
		if p.current == tokens.Identifier && p.pseudo != tokens.InvalidPseudoToken {
			p.report("Expected %s, received %s", t.String(), p.pseudo.String())
		} else {
			p.report("Expected %s, received %v", t.String(), p.current.String())
		}
	}
}

func (p *parser) expectPseudoSymbol(t tokens.PseudoToken) {
	if p.pseudo == t {
		p.next()
	} else {
		p.builder.PushContext()
		defer p.builder.PopContext()
		if p.current == tokens.Symbol && p.pseudo != tokens.InvalidPseudoToken {
			p.report("Expected %s, received %s", t.String(), p.pseudo.String())
		} else {
			p.report("Expected %s, received %s", t.String(), p.current.String())
		}
	}
}

func (p *parser) expects(ts ...tokens.Token) ast.Element {
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
	tokens.LBrace, tokens.LParen, tokens.Let,
}

func (p *parser) expression() ast.Element {
	switch p.current {
	case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
		tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LBrace,
		tokens.LParen, tokens.Let:
		return p.simpleExpression()
	default:
		return p.expects(primitiveTokens...)
	}
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
		tokens.LBrace, tokens.LParen:
		for {
			switch p.current {
			case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
				tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
				tokens.LBrace, tokens.LParen:
				argument := p.argument()
				arguments = append(arguments, argument)
				if p.current == tokens.Comma {
					p.next()
					switch p.current {
					case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
						tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier,
						tokens.LBrace, tokens.LParen:
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
	p.expectPseudoSymbol(tokens.Equal)
	value := p.expression()
	return p.builder.NamedArgument(name, value)
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
	p.expectPseudoSymbol(tokens.Bar)
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
	p.expectPseudoSymbol(tokens.Equal)
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
	p.expectPseudoSymbol(tokens.Arrow)
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
		p.expectPseudoSymbol(tokens.Equal)
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
	p.expectPseudoSymbol(tokens.Spread)
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
