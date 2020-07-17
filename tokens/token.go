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

	// Add '+'
	Add

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

	// LParen '('
	LParen

	// RParen ')'
	RParen

	// LBrack '['
	LBrack

	// RBrack ']'
	RBrack

	// LBrace '{'
	LBrace

	// RBrace '}'
	RBrace

	// Semi ';'
	Semi

	// Colon ':'
	Colon

	// Question '?'
	Question

	// Not '!'
	Not

	// Comma ','
	Comma

	// Dot '.'
	Dot

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

	// NotEqual '!='
	NotEqual

	// LessThan '<'
	LessThan

	// LessThanEqual '<='
	LessThanEqual

	// Arrow '->'
	Arrow

	// Range '..'
	Range

	// Assign ':='
	Assign

	// Platform '!!'
	Platform

	// As 'as'
	As

	// Boolean 'boolean'
	Boolean

	// Break 'break'
	Break

	// Case 'case'
	Case

	// Continue 'continue'
	Continue

	// Constraint 'constraint'
	Constraint

	// Default 'default'
	Default

	// Else 'else'
	Else

	// Enum 'enum'
	Enum

	// False 'false'
	False

	// Float 'float'
	Float

	// For 'for'
	For

	// If 'if'
	If

	// In 'in'
	In

	// Int 'int'
	Int

	// Interface 'interface
	Interface

	// Let 'let'
	Let

	// Match 'match'
	Match

	// Method 'method'
	Method

	// Operator 'operator'
	Operator

	// Property 'property'
	Property

	// Record 'record'
	Record

	// Return 'return'
	Return

	// String 'string'
	String

	// Switch 'switch'
	Switch

	// True 'true'
	True

	// Type 'type'
	Type

	// Value 'value'
	Value

	// Var 'var'
	Var

	// Void 'void'
	Void

        // Where 'where'
        Where

	lastToken
)

var tokens = [...]string{
	Invalid:          "<invalid>",
	EOF:              "<eof>",
	Identifier:       "<identifier>",
	LiteralByte:      "<byte>",
	LiteralInt:       "<int>",
	LiteralUInt:      "<uint>",
	LiteralLong:      "<long>",
	LiteralULong:     "<ulong>",
	LiteralFloat:     "<float>",
	LiteralDouble:    "<double>",
	LiteralString:    "<string>",
	LiteralRune:      "<rune>",
	Add:              "+",
        Bar:              "|",
	Sub:              "-",
	Mult:             "*",
	Div:              "/",
	Rem:              "%",
	LParen:           "(",
	RParen:           ")",
	LBrack:           "[",
	RBrack:           "]",
	LBrace:           "{",
	RBrace:           "}",
	Semi:             ";",
	Colon:            ":",
	Question:         "?",
	Not:              "!",
	Comma:            ",",
	Dot:              ".",
	LogicalAnd:       "&&",
	LogicalOr:        "||",
	GreaterThan:      ">",
	GreaterThanEqual: ">=",
	Equal:            "=",
	NotEqual:         "!=",
	LessThan:         "<",
	LessThanEqual:    "<=",
	Arrow:            "->",
	Range:            "..",
	Assign:           ":=",
	Platform:         "!!",
	As:               "as",
	Boolean:          "boolean",
	Break:            "break",
	Case:             "case",
	Constraint:       "constraint",
	Continue:         "continue",
	Default:          "default",
	Else:             "else",
	Enum:             "enum",
	False:            "false",
	Float:            "float",
	For:              "for",
	If:               "if",
	In:               "in",
	Int:              "int",
	Interface:        "interface",
	Let:              "let",
	Match:            "match",
	Method:           "method",
	Operator:         "operator",
	Property:         "property",
	Record:           "record",
	Return:           "return",
	String:           "string",
	Switch:           "switch",
	True:             "true",
	Type:             "type",
	Value:            "value",
	Var:              "var",
	Void:             "void",
        Where:            "where",
}

func (t Token) String() string {
	if t >= 0 && t < lastToken {
		return tokens[t]
	}
	return "<unknown>"
}
