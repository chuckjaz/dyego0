package parser

import (
    "fmt"

    "dyego0/ast"
    "dyego0/tokens"
)

// Parser parses text and returns an ast element
type Parser interface {
    Errors() []ast.Error
    Parse() ast.Element
}

type parser struct {
    scanner *Scanner
    builder ast.Builder
    current tokens.Token
    errors []ast.Error
}

// NewParser creates a new parser
func NewParser(scanner *Scanner) Parser {
    builder := ast.NewBuilder(scanner)
    p := &parser{scanner: scanner, builder: builder}
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
        p.report("Expected %v, recieved %v", t, p.current)
        p.next()
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

func (p *parser) next() tokens.Token {
    p.builder.UpdateContext()
    p.current = p.scanner.Next()
    return p.current
}

func (p *parser) expectIdent() ast.Name {
    if p.current == tokens.Identifier {
        result := p.builder.Name(p.scanner.Value().(string))
        p.next()
        return result
    }
    result := p.builder.Name("<error>")
    p.expect(tokens.Identifier)
    return result
}

var primitiveTokens = []tokens.Token{
    tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
        tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LParen,
}

func (p* parser) expression() ast.Element {
    switch p.current {
    case tokens.LiteralRune, tokens.LiteralByte, tokens.LiteralInt, tokens.LiteralUInt, tokens.LiteralLong, tokens.LiteralULong,
        tokens.LiteralDouble, tokens.LiteralFloat, tokens.LiteralString, tokens.True, tokens.False, tokens.Identifier, tokens.LParen:
        return p.primitive()
    default:
        return p.expects(primitiveTokens...)
    }
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
    case tokens.LParen:
        expr := p.expression()
        p.expect(tokens.RParen)
        return expr
    default:
        return p.expects(primitiveTokens...)
    }
}
