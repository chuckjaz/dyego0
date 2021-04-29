package tokens

// Token is a lexical token in a file
type Token int

const (
	// Invalid token
	Invalid Token = iota

	// EOF is the end of the file
	EOF

	// Identifier is an identifier (e.g. SomeIdent)
	Identifier

	// Symbol is a string of special symbols typecially used as operators (e.g. +, -, !=, etc.)
	Symbol

	// LiteralString is a string literal (e.g. "string")
	LiteralString

	// LiteralRune is a rune literal (e.g. 'b')
	LiteralRune

	// LiteralInt is a literal integer (e.g. 12345)
	LiteralInt

	// LiteralByte is a literal byte (e.g. 27b)
	LiteralByte

	// LiteralUInt is a literal unsigned int (e.g. 23u)
	LiteralUInt

	// LiteralLong is a literal long (e.g. 122233322l)
	LiteralLong

	// LiteralULong is a literal unsigned long (e.g. 12332323ul)
	LiteralULong

	// LiteralFloat is a literal float (e.g. 123.5f)
	LiteralFloat

	// LiteralDouble is a literal double (e.g. 123.5)
	LiteralDouble

	// LParen '('
	LParen

	// RParen ')'
	RParen

	// LBrack '['
	LBrack

	// RBrack ']'
	RBrack

	// LBrackBang '[!'
	LBrackBang

	// BangRBrack '!]'
	BangRBrack

	// LBrace '{'
	LBrace

	// RBrace '}'
	RBrace

	// LBraceBang '{!'
	LBraceBang

	// BangRBrace '!}'
	BangRBrace

	// Semi ';'
	Semi

	// Colon ':'
	Colon

	// Comma ','
	Comma

	// Dot '.'
	Dot

	// Scope '::'
	Scope

	// VocabularyStart "<|"
	VocabularyStart

	// VocabularyEnd "|>"
	VocabularyEnd

	// False 'false'
	False

	// Let 'let'
	Let

	// True 'true'
	True

	// Return 'return'
	Return

	// Val 'val'
	Val

	// Var 'var'
	Var

	lastToken

	// InvalidToken is an out of band token value. Not all invalid tokens equal InvalidToken but InvalidToken
	// is guaranteed to be invalid.
	InvalidToken
)

var tokens = [...]string{
	Invalid:         "<invalid>",
	EOF:             "<eof>",
	Identifier:      "<identifier>",
	Symbol:          "<symbol>",
	LiteralByte:     "<byte>",
	LiteralInt:      "<int>",
	LiteralUInt:     "<uint>",
	LiteralLong:     "<long>",
	LiteralULong:    "<ulong>",
	LiteralFloat:    "<float>",
	LiteralDouble:   "<double>",
	LiteralString:   "<string>",
	LiteralRune:     "<rune>",
	LParen:          "(",
	RParen:          ")",
	LBrack:          "[",
	RBrack:          "]",
	LBrackBang:      "[!",
	BangRBrack:      "!]",
	LBrace:          "{",
	RBrace:          "}",
	LBraceBang:      "{!",
	BangRBrace:      "!}",
	Semi:            ";",
	Colon:           ":",
	Comma:           ",",
	Dot:             ".",
	Scope:           "::",
	VocabularyStart: "<|",
	VocabularyEnd:   "|>",
	False:           "false",
	Let:             "let",
	True:            "true",
	Return:          "return",
	Val:             "val",
	Var:             "var",
}

func (t Token) String() string {
	if t >= 0 && t < lastToken {
		return tokens[t]
	}
	return "<invalid>"
}

const ()

// PseudoToken are special symbols and identifier that are recoginized
// by the scanner but reported as an indentifier or symbol token
type PseudoToken int

const (
	// After is the pseudo token "after"
	After PseudoToken = iota

	// Before is the psuedo token "before"
	Before

	// Break is a pseudo token "break"
	Break

	// Continue is a pseudo token "continue"
	Continue

	// Else is a pseudo token "else"
	Else

	// Identifiers is a pseudo token "identifiers"
	Identifiers

	// If is a pseudo token "if"
	If

	// Infix is the pseudo token "infix"
	Infix

	// Left is the pseudo token "left"
	Left

	// Loop is the pseudo token "loop"
	Loop

	// Operator it he pseudo token "opeartor"
	Operator

	// Postfix is the pseudo token "postfix"
	Postfix

	// Prefix is the pseudo token "prefix"
	Prefix

	// Right is the pseudo token "right"
	Right

	// When is the pseudo token "when"
	When

	// Where is a pseudo token "where"
	Where

	// Add '+'
	Add

	// And '&'
	And

	// Bar '|'
	Bar

	// Sub '-'
	Sub

	// Mult '*'
	Mult

	// Div '/'
	Div

	// Rem '%'
	Rem

	// Not '!'
	Not

	// LogicalAnd '&&'
	LogicalAnd

	// LogicalOr '||'
	LogicalOr

	// GreaterThan '>'
	GreaterThan

	// GreaterThanEqual '>='
	GreaterThanEqual

	// Equal '='
	Equal

	// DoubleEqual '=='
	DoubleEqual

	// NotEqual '!='
	NotEqual

	// LessThan '<'
	LessThan

	// LessThanEqual '<='
	LessThanEqual

	// Question '?'
	Question

	// Arrow '->'
	Arrow

	// Range '..'
	Range

	// Spread '...''
	Spread

	lastPseudoToken

	// InvalidPseudoToken is an out of band value for pseudo token
	InvalidPseudoToken

	// Escaped is an escaped identifier using the `...` syntax
	Escaped
)

var pseudoTokens = [...]string{
	After:            "after",
	Before:           "before",
	Break:            "break",
	Continue:         "continue",
	Else:             "else",
	Identifiers:      "identifiers",
	If:               "if",
	Infix:            "infix",
	Left:             "left",
	Loop:             "loop",
	Operator:         "operator",
	Prefix:           "prefix",
	Postfix:          "postfix",
	Right:            "right",
	When:             "when",
	Where:            "where",
	Add:              "+",
	And:              "&",
	Bar:              "|",
	Sub:              "-",
	Mult:             "*",
	Div:              "/",
	Rem:              "%",
	Not:              "!",
	LogicalAnd:       "&&",
	LogicalOr:        "||",
	GreaterThan:      ">",
	GreaterThanEqual: ">=",
	Equal:            "=",
	DoubleEqual:      "==",
	NotEqual:         "!=",
	LessThan:         "<",
	LessThanEqual:    "<=",
	Question:         "?",
	Arrow:            "->",
	Range:            "..",
	Spread:           "...",
	Escaped:          "<escaped>",
}

func (t PseudoToken) String() string {
	if t >= 0 && t < lastPseudoToken {
		return pseudoTokens[t]
	}
	return "<invalid>"
}
