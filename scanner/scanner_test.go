package scanner

import (
	"go/token"
	"strings"
	"testing"

	"dyego0/tokens"
)

func paniced(fn func()) bool {
	paniced := false
	func() {
		defer func() {
			paniced = recover() != nil
		}()
		fn()
	}()
	return paniced
}

func TestNew(t *testing.T) {
	NewScanner([]byte{0}, 0)
}

func TestInvalidBuffer(t *testing.T) {
	emptyInvalid := paniced(func() {
		NewScanner([]byte{}, 0)
	})
	noNullTerminatorInvalid := paniced(func() {
		NewScanner([]byte{'a', 'b', 'c'}, 0)
	})
	if !emptyInvalid || !noNullTerminatorInvalid {
		t.Fail()
	}
}

func parsePseudoBytes(t *testing.T, src []byte, expectedToken tokens.Token, expected ...tokens.PseudoToken) (token.Pos, int, int) {
	scanner := NewScanner(src, 0)
	var received tokens.PseudoToken
	for _, pseudoToken := range expected {
		token := scanner.Next()
		if token != expectedToken {
			t.Errorf("Expected '%s', received '%s'", expectedToken.String(), token.String())
		}
		received = scanner.PseudoToken()
		if pseudoToken != received {
			t.Errorf("Expected '%s', received '%s'", pseudoToken.String(), received.String())
		}
		if pseudoToken != tokens.InvalidPseudoToken && pseudoToken.String() != scanner.Value() {
			t.Errorf("Expected pseudoToken text to match string text %s, %v", pseudoToken, scanner.Value())
		}
	}
	if scanner.Next() != tokens.EOF {
		t.Error("Not enough tokens")
	}
	return scanner.Start(), scanner.Offset(), scanner.Line()
}

func parsePseudo(t *testing.T, text string, expectedToken tokens.Token, expected ...tokens.PseudoToken) (token.Pos, int, int) {
	return parsePseudoBytes(t, append([]byte(text), 0), expectedToken, expected...)
}

func parseBytes(t *testing.T, src []byte, expected ...tokens.Token) (token.Pos, int, int) {
	scanner := NewScanner(src, 0)
	var received tokens.Token
	for _, token := range expected {
		received = scanner.Next()
		if token != received {
			t.Errorf("Expected '%v', received '%v'", token, received)
		}
	}
	if scanner.Next() != tokens.EOF {
		t.Error("Not enough tokens")
	}
	return scanner.Start(), scanner.Offset(), scanner.Line()
}

func parseString(t *testing.T, text string, expected ...tokens.Token) (token.Pos, int, int) {
	src := append([]byte(text), 0)
	return parseBytes(t, src, expected...)
}

func scanOne(text string) (scanner *Scanner, token tokens.Token) {
	src := append([]byte(text), 0)
	scanner = NewScanner(src, 0)
	token = scanner.Next()
	return
}

func rejectOne(t *testing.T, text string, msg string) {
	scanner, token := scanOne(text)
	if token != tokens.Invalid {
		t.Error("Expected an invalid token")
	}
	if !strings.Contains(scanner.Message(), msg) {
		t.Errorf("Expected %s to contain %s", scanner.Message(), msg)
	}
}

func TestIdentifer(t *testing.T) {
	parseString(t, "ident", tokens.Identifier)
	parseString(t, "  ident   ident2 _ _12", tokens.Identifier, tokens.Identifier, tokens.Identifier, tokens.Identifier)
}

func TestReservedSymbols(t *testing.T) {
	parseString(t, "{}()[];: ,.::<||>[!!]",
		tokens.LBrace, tokens.RBrace, tokens.LParen, tokens.RParen, tokens.LBrack, tokens.RBrack,
		tokens.Semi, tokens.Colon, tokens.Comma, tokens.Dot, tokens.Scope, tokens.VocabularyStart,
		tokens.VocabularyEnd, tokens.LBrackBang, tokens.BangRBrack,
	)
}

func TestPseudoSymbols(t *testing.T) {
	parsePseudo(t, "+ & | - * / % ! && || > >= = == != < <= ? -> .. ...", tokens.Symbol,
		tokens.Add, tokens.And, tokens.Bar, tokens.Sub, tokens.Mult, tokens.Div, tokens.Rem, tokens.Not,
		tokens.LogicalAnd, tokens.LogicalOr, tokens.GreaterThan, tokens.GreaterThanEqual,
		tokens.Equal, tokens.DoubleEqual, tokens.NotEqual, tokens.LessThan, tokens.LessThanEqual,
		tokens.Question, tokens.Arrow, tokens.Range, tokens.Spread,
	)
}

func testReservedWord(t *testing.T, tokenList ...tokens.Token) {
	for _, token := range tokenList {
		text := token.String()
		example := text + " " + text + " " + text + "r " + " " + text
		parseString(t, example, token, token, tokens.Identifier, token)
	}
}

func testPseudo(t *testing.T, expectedToken tokens.Token, extender string, pseudoList ...tokens.PseudoToken) {
	for _, pseudo := range pseudoList {
		text := pseudo.String()
		example := text + " " + text + " " + text + extender + " " + text
		parsePseudo(t, example, expectedToken, pseudo, pseudo, tokens.InvalidPseudoToken, pseudo)
	}
}

func testPseudoWord(t *testing.T, pseudoList ...tokens.PseudoToken) {
	testPseudo(t, tokens.Identifier, "r", pseudoList...)
}

func testPseudoSymbol(t *testing.T, pseudoList ...tokens.PseudoToken) {
	testPseudo(t, tokens.Symbol, "+", pseudoList...)
}

func TestReservedWords(t *testing.T) {
	testReservedWord(t, tokens.False, tokens.Let, tokens.True, tokens.Return, tokens.Val, tokens.Var)
}

func TestPseudoReservedSymbols(t *testing.T) {
	// Test pseudo symbols that obey extenders rule
	testPseudoSymbol(t,
		tokens.Add, tokens.Bar, tokens.Sub, tokens.Mult, tokens.Div, tokens.Rem, tokens.Not,
		tokens.LogicalAnd, tokens.LogicalOr, tokens.GreaterThan, tokens.GreaterThanEqual,
		tokens.Equal, tokens.DoubleEqual, tokens.NotEqual, tokens.LessThan, tokens.LessThanEqual,
	)
}

func TestPseudoWords(t *testing.T) {
	testPseudoWord(t,
		tokens.After, tokens.Before, tokens.Break, tokens.Continue, tokens.Else, tokens.Infix, tokens.Identifiers,
		tokens.Left, tokens.Loop, tokens.Operator, tokens.Postfix, tokens.Prefix, tokens.Right, tokens.When, tokens.Where,
	)
}

func TestLineCount(t *testing.T) {
	_, _, lines := parseString(t, " \n  \r\n  \n ident ", tokens.Identifier)
	if lines != 4 {
		t.Errorf("Expected %d to be 4", lines)
	}
}

func TestOffset(t *testing.T) {
	src := " a b c "
	_, offset, _ := parseString(t, src, tokens.Identifier, tokens.Identifier, tokens.Identifier)
	if offset != len(src) {
		t.Errorf("Expected %d to be %d", offset, len(src))
	}
}

func TestStart(t *testing.T) {
	text := " a b c "
	src := append([]byte(text), 0)
	scanner := NewScanner(src, 0)
	scanner.Next()
	for _, expected := range []token.Pos{1, 3, 5} {
		if scanner.Start() != expected {
			t.Errorf("Expected %d, to be %d", scanner.Start(), expected)
		}
		scanner.Next()
	}
}

func TestInt(t *testing.T) {
	scanner, token := scanOne(" 10 ")
	if token != tokens.LiteralInt {
		t.Error("Expected int literal")
	}
	if scanner.Value().(int) != 10 {
		t.Error("Expected 10")
	}
}

func TestLong(t *testing.T) {
	scanner, token := scanOne(" 10l ")
	if token != tokens.LiteralLong {
		t.Error("Expected long literal")
	}
	if scanner.Value().(int64) != 10 {
		t.Error("Expected 10")
	}
}

func TestUInt(t *testing.T) {
	scanner, token := scanOne(" 10u ")
	if token != tokens.LiteralUInt {
		t.Error("Expected uint literal")
	}
	if scanner.Value().(uint) != 10 {
		t.Error("Expected 10")
	}
}

func TestULInt(t *testing.T) {
	scanner, token := scanOne(" 10ul ")
	if token != tokens.LiteralULong {
		t.Error("Expected ulong literal")
	}
	if scanner.Value().(uint64) != 10 {
		t.Error("Expected 10")
	}
}

func TestByte(t *testing.T) {
	scanner, token := scanOne(" 10b ")
	if token != tokens.LiteralByte {
		t.Error("Expected integer byte")
	}
	if scanner.Value().(byte) != byte(10) {
		t.Error("Expected 10")
	}
}

func TestFloat(t *testing.T) {
	scanner, token := scanOne(" 10.1f ")
	if token != tokens.LiteralFloat {
		t.Error("Expected float literal")
	}
	if scanner.Value().(float32) != float32(10.1) {
		t.Error("Expected 10.1f")
	}
}

func TestDouble(t *testing.T) {
	scanner, token := scanOne(" 10.1 ")
	if token != tokens.LiteralDouble {
		t.Error("Expected double literal")
	}
	if scanner.Value().(float64) != 10.1 {
		t.Error("Expected 10.1")
	}
}

func TestString(t *testing.T) {
	scanner, token := scanOne(" \"this is a test\" ")
	if token != tokens.LiteralString {
		t.Error("Expected string literal")
	}
	if scanner.Value().(string) != "this is a test" {
		t.Error("Expected \"this is a test\"")
	}
}

func TestStringEscapes(t *testing.T) {
	scanner, token := scanOne(" \"[[ \\n \\b \\t \\\" ]]\" ")
	if token != tokens.LiteralString {
		t.Error("Expected string literal")
	}
	var received = scanner.Value().(string)
	var expected = "[[ \n \b \t \" ]]"
	if received != expected {
		t.Errorf("Expected %s to equal %s", received, expected)
	}
}

func testRuneValue(t *testing.T, text string, r rune) {
	scanner, token := scanOne(text)
	if token != tokens.LiteralRune {
		t.Error("Expected rune literal")
	}
	var received = scanner.Value().(rune)
	if received != r {
		t.Errorf("Expected %c to equal %c", received, r)
	}
}

func TestRune(t *testing.T) {
	testRuneValue(t, " 'a'  ", 'a')
	testRuneValue(t, " '\n'  ", '\n')
}

func TestInvalidFloat(t *testing.T) {
	rejectOne(t, "10.1.1", "invalid syntax")
}

func TestInvalidString(t *testing.T) {
	rejectOne(t, " \"Test \n \"", "Unterminated string")
}
