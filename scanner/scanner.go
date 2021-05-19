package scanner

import (
	"fmt"
	"strconv"

	"dyego0/location"
	"dyego0/tokens"
)

const (
	// InternalScan enables internal identifiers
	InternalScan = 1 << iota
)

// Scanner is a Dyego scanner
type Scanner struct {
	src    []byte
	fb     tokens.FileBuilder
	offset int
	line   int
	start  int
	end    int
	nlloc  int
	msg    string
	flags  int
	pseudo tokens.PseudoToken
	value  interface{}
}

// NewScanner creates a scanner
func NewScanner(src []byte, flags int, fb tokens.FileBuilder) *Scanner {
	length := len(src)
	if length == 0 || src[length-1] != 0 {
		panic("NewScanner: src must be null terminated")
	}
	if fb != nil {
		fb.AddLine(0)
	}
	return &Scanner{src: src, fb: fb, line: 1, nlloc: -1, flags: flags}
}

// Clone preserves a copy of the scanner at the current state which can then be used
// for backtracking, if necessary by using the returned instance instead of the
// instance that was moved forward.
func (s *Scanner) Clone() *Scanner {
	return &Scanner{src: s.src, fb: s.fb, offset: s.offset, line: s.line, start: s.start, end: s.end, nlloc: s.nlloc,
		msg: s.msg, flags: s.flags, pseudo: s.pseudo, value: s.value}
}

// Line is the current line of the scanner
func (s *Scanner) Line() int {
	return s.line
}

// Start is the start of the current token
func (s *Scanner) Start() location.Pos {
	if s.fb != nil {
		return s.fb.Pos(s.start)
	}
	return location.Pos(s.start)
}

// End is the end of the current token
func (s *Scanner) End() location.Pos {
	if s.fb != nil {
		return s.fb.Pos(s.end)
	}
	return location.Pos(s.end)
}

// NewLineLocation is the location of a new line prior to the current token
func (s *Scanner) NewLineLocation() location.Pos {
	if s.nlloc >= 0 && s.fb != nil {
		return s.fb.Pos(s.nlloc)
	}
	return location.Pos(s.nlloc)
}

// Offset is the current location of the scanner
func (s *Scanner) Offset() int {
	return s.offset
}

// Value is the value of the last literal
func (s *Scanner) Value() interface{} {
	return s.value
}

// Message is the error message if there is one
func (s *Scanner) Message() string {
	return s.msg
}

func identExtender(b byte) bool {
	switch b {
	case 'a', 'b', 'c', 'd', 'e',
		'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E',
		'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z',
		'_', '$',
		'0', '1', '2', '3', '4',
		'5', '6', '7', '8', '9':
		return true
	}
	return false
}

func symbolExtender(b byte) bool {
	switch b {
	case '~', '!', '@', '#', '$', '%', '^', '&',
		'-', '_', '+', '=', '/', '?', '|', ':':
		return true
	}
	return false
}

// PseudoToken returns the pseudo token for the current identifier
func (s *Scanner) PseudoToken() tokens.PseudoToken {
	return s.pseudo
}

// Next moves the scanner to the next token
func (s *Scanner) Next() tokens.Token {
	s.end = s.offset
	offset := s.offset
	start := s.offset
	src := s.src
	line := s.line
	result := tokens.Invalid
	s.pseudo = tokens.InvalidPseudoToken
	s.value = nil
	s.nlloc = -1
loop:
	for {
		b := src[offset]
		start = offset
		offset++
		switch b {
		// EOF
		case 0:
			result = tokens.EOF
			offset--
			break loop

			// Whitespace
		case ' ', '\t':
			continue
		case '\r':
			if src[offset] == '\n' {
				offset++
			}
			fallthrough
		case '\n':
			line++
			s.nlloc = offset - 1
			if s.fb != nil {
				s.fb.AddLine(offset)
			}
			continue

		case '+', '|', '-', '*', '/', '%', '!', '&', '>', '<',
			'=', ':', '?':
			// Pseudo-symbols
			switch b {
			case '+':
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.Add
					s.value = "+"
					break loop
				}
			case '|':
				switch src[offset] {
				case '>':
					offset++
					result = tokens.VocabularyEnd
					break loop
				case '|':
					if !symbolExtender(src[offset+1]) {
						offset++
						result = tokens.Symbol
						s.pseudo = tokens.LogicalOr
						s.value = "||"
						break loop
					}
				}
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.Bar
					s.value = "|"
					break loop
				}
			case '-':
				if src[offset] == '>' && !symbolExtender(src[offset+1]) {
					offset++
					result = tokens.Symbol
					s.pseudo = tokens.Arrow
					s.value = "->"
					break loop
				}
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.Sub
					s.value = "-"
					break loop
				}
			case '*':
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.Mult
					s.value = "*"
					break loop
				}
			case '/':
				if src[offset] == '/' {
				commentLoop:
					for {
						b := src[offset]
						offset++
						switch b {
						case '\n', '\r', 0:
							offset--
							break commentLoop
						}
					}
					continue loop
				}
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.Div
					s.value = "/"
					break loop
				}
			case '%':
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.Rem
					s.value = "%"
					break loop
				}
			case '!':
				switch src[offset] {
				case '=':
					if !symbolExtender(src[offset+1]) {
						offset++
						result = tokens.Symbol
						s.pseudo = tokens.NotEqual
						s.value = "!="
						break loop
					}
				case ']':
					offset++
					result = tokens.BangRBrack
					s.value = "!]"
					break loop
				case '}':
					offset++
					result = tokens.BangRBrace
					s.value = "!}"
					break loop
				default:
					if !symbolExtender(src[offset]) {
						result = tokens.Symbol
						s.pseudo = tokens.Not
						s.value = "!"
						break loop
					}
				}
			case '&':
				if src[offset] == '&' && !symbolExtender(src[offset+1]) {
					offset++
					result = tokens.Symbol
					s.pseudo = tokens.LogicalAnd
					s.value = "&&"
					break loop
				}
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.And
					s.value = "&"
					break loop
				}
			case '<':
				switch src[offset] {
				case '=':
					if !symbolExtender(src[offset+1]) {
						offset++
						result = tokens.Symbol
						s.pseudo = tokens.LessThanEqual
						s.value = "<="
						break loop
					}
				case '|':
					offset++
					result = tokens.VocabularyStart
					break loop
				}
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.LessThan
					s.value = "<"
					break loop
				}
			case '>':
				if src[offset] == '=' && !symbolExtender(src[offset+1]) {
					offset++
					result = tokens.Symbol
					s.pseudo = tokens.GreaterThanEqual
					s.value = ">="
					break loop
				}
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.GreaterThan
					s.value = ">"
					break loop
				}
			case '=':
				if src[offset] == '=' && !symbolExtender(src[offset+1]) {
					offset++
					result = tokens.Symbol
					s.pseudo = tokens.DoubleEqual
					s.value = "=="
					break loop
				}
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.Equal
					s.value = "="
					break loop
				}
			case ':':
				if src[offset] == ':' && !symbolExtender(src[offset+1]) {
					offset++
					result = tokens.Scope
					break loop
				}
				if !symbolExtender(src[offset]) {
					result = tokens.Colon
					break loop
				}
			case '?':
				if !symbolExtender(src[offset]) {
					result = tokens.Symbol
					s.pseudo = tokens.Question
					s.value = "?"
					break loop
				}
			}
			fallthrough
		case '~', '@', '#', '$', '^':
		symbolLoop:
			for {
				last := offset
				b = src[offset]
				offset++
				switch b {
				default:
					offset = last
					break symbolLoop
				case '~', '!', '@', '#', '$', '%', '^', '&',
					'-', '_', '+', '=', ',', '/', '?',
					'|', ':':
					continue
				}
			}
			s.value = string(src[start:offset])
			result = tokens.Symbol
		case ',':
			result = tokens.Comma
			s.value = ","
		case ';':
			result = tokens.Semi
			s.value = ";"
		case '{':
			if src[offset] == '!' {
				offset++
				result = tokens.LBraceBang
				s.value = "{!"
			} else {
				result = tokens.LBrace
				s.value = "{"
			}
		case '}':
			result = tokens.RBrace
			s.value = "}"
		case '[':
			if src[offset] == '!' {
				offset++
				result = tokens.LBrackBang
				s.value = "[!"
			} else {
				result = tokens.LBrack
				s.value = "["
			}
		case ']':
			result = tokens.RBrack
			s.value = "]"
		case '(':
			result = tokens.LParen
			s.value = "("
		case ')':
			result = tokens.RParen
			s.value = ")"
		case '.':
			if src[offset] == '.' {
				if src[offset+1] == '.' {
					offset += 2
					result = tokens.Symbol
					s.pseudo = tokens.Spread
					s.value = "..."
				} else {
					offset++
					result = tokens.Symbol
					s.pseudo = tokens.Range
					s.value = ".."
				}
			} else {
				result = tokens.Dot
				s.value = "."
			}

		// Pseudo reserved words and identifiers
		case 'a', 'b', 'c', 'e', 'f', 'i', 'l', 'o', 'p', 't', 'r', 'v', 'w':
			switch b {
			case 'a':
				switch src[offset] {
				case 'f':
					// after
					if src[offset+1] == 't' && src[offset+2] == 'e' && src[offset+3] == 'r' &&
						!identExtender(src[offset+4]) {
						offset += 4
						result = tokens.Identifier
						s.pseudo = tokens.After
						s.value = "after"
						break loop
					}
				}
			case 'b':
				switch src[offset] {
				case 'e':
					// before
					if src[offset+1] == 'f' && src[offset+2] == 'o' && src[offset+3] == 'r' && src[offset+4] == 'e' &&
						!identExtender(src[offset+5]) {
						offset += 5
						result = tokens.Identifier
						s.pseudo = tokens.Before
						s.value = "before"
						break loop
					}
				case 'r':
					// break
					if src[offset+1] == 'e' && src[offset+2] == 'a' && src[offset+3] == 'k' && !identExtender(src[offset+4]) {
						offset += 4
						result = tokens.Identifier
						s.pseudo = tokens.Break
						s.value = "break"
						break loop
					}
				}
			case 'c':
				// continue
				if src[offset] == 'o' && src[offset+1] == 'n' && src[offset+2] == 't' && src[offset+3] == 'i' &&
					src[offset+4] == 'n' && src[offset+5] == 'u' && src[offset+6] == 'e' && !identExtender(src[offset+7]) {
					offset += 7
					result = tokens.Identifier
					s.pseudo = tokens.Continue
					s.value = "continue"
					break loop
				}
			case 'e':
				// else
				if src[offset] == 'l' && src[offset+1] == 's' && src[offset+2] == 'e' && !identExtender(src[offset+3]) {
					offset += 3
					result = tokens.Identifier
					s.pseudo = tokens.Else
					s.value = "else"
					break loop
				}
			case 'f':
				// false
				if src[offset] == 'a' && src[offset+1] == 'l' && src[offset+2] == 's' && src[offset+3] == 'e' &&
					!identExtender(src[offset+4]) {
					offset += 4
					result = tokens.False
					s.value = "false"
					break loop
				}
			case 'i':
				switch src[offset] {
				case 'd':
					// identifiers
					if src[offset+1] == 'e' && src[offset+2] == 'n' && src[offset+3] == 't' && src[offset+4] == 'i' &&
						src[offset+5] == 'f' && src[offset+6] == 'i' && src[offset+7] == 'e' && src[offset+8] == 'r' &&
						src[offset+9] == 's' && !identExtender(src[offset+10]) {
						offset += 10
						result = tokens.Identifier
						s.pseudo = tokens.Identifiers
						s.value = "identifiers"
						break loop
					}
				case 'f':
					if !identExtender(src[offset+1]) {
						offset++
						result = tokens.Identifier
						s.pseudo = tokens.If
						s.value = "if"
						break loop
					}
				case 'n':
					switch src[offset+1] {
					case 'f':
						// infix
						if src[offset+2] == 'i' && src[offset+3] == 'x' && !identExtender(src[offset+4]) {
							offset += 4
							result = tokens.Identifier
							s.pseudo = tokens.Infix
							s.value = "infix"
							break loop
						}
					}
				}
			case 'l':
				switch src[offset] {
				case 'e':
					switch src[offset+1] {
					case 'f':
						// left
						if src[offset+2] == 't' && !identExtender(src[offset+3]) {
							offset += 3
							result = tokens.Identifier
							s.pseudo = tokens.Left
							s.value = "left"
							break loop
						}
					case 't':
						// let
						if !identExtender(src[offset+2]) {
							offset += 2
							result = tokens.Let
							break loop
						}
					}
				case 'o':
					// loop
					if src[offset+1] == 'o' && src[offset+2] == 'p' && !identExtender(src[offset+3]) {
						offset += 3
						result = tokens.Identifier
						s.pseudo = tokens.Loop
						s.value = "loop"
						break loop
					}
				}
			case 'o':
				// operator
				if src[offset] == 'p' && src[offset+1] == 'e' && src[offset+2] == 'r' &&
					src[offset+3] == 'a' && src[offset+4] == 't' && src[offset+5] == 'o' &&
					src[offset+6] == 'r' && !identExtender(src[offset+7]) {
					offset += 7
					result = tokens.Identifier
					s.pseudo = tokens.Operator
					s.value = "operator"
					break loop
				}
			case 'p':
				switch src[offset] {
				case 'o':
					// postfix
					if src[offset+1] == 's' && src[offset+2] == 't' && src[offset+3] == 'f' &&
						src[offset+4] == 'i' && src[offset+5] == 'x' && !identExtender(src[offset+6]) {
						offset += 6
						result = tokens.Identifier
						s.pseudo = tokens.Postfix
						s.value = "postfix"
						break loop
					}
				case 'r':
					switch src[offset+1] {
					case 'e':
						// prefix
						if src[offset+2] == 'f' && src[offset+3] == 'i' && src[offset+4] == 'x' &&
							!identExtender(src[offset+5]) {
							offset += 5
							result = tokens.Identifier
							s.pseudo = tokens.Prefix
							s.value = "prefix"
							break loop
						}
					}
				}
			case 't':
				// true
				if src[offset] == 'r' && src[offset+1] == 'u' && src[offset+2] == 'e' &&
					!identExtender(src[offset+3]) {
					offset += 3
					result = tokens.True
					s.value = "true"
					break loop
				}
			case 'r':
				switch src[offset] {
				case 'e':
					switch src[offset+1] {
					case 't':
						// return
						if src[offset+2] == 'u' && src[offset+3] == 'r' && src[offset+4] == 'n' &&
							!identExtender(src[offset+5]) {
							offset += 5
							result = tokens.Return
							break loop
						}
					}
				case 'i':
					// right
					if src[offset+1] == 'g' && src[offset+2] == 'h' && src[offset+3] == 't' &&
						!identExtender(src[offset+4]) {
						offset += 4
						result = tokens.Identifier
						s.pseudo = tokens.Right
						s.value = "right"
						break loop
					}
				}
			case 'v':
				switch src[offset] {
				case 'a':
					switch src[offset+1] {
					case 'r':
						// var
						if !identExtender(src[offset+2]) {
							offset += 2
							result = tokens.Var
							break loop
						}
					case 'l':
						// val
						if !identExtender(src[offset+2]) {
							offset += 2
							result = tokens.Val
							break loop
						}
					}
				}
			case 'w':
				switch src[offset] {
				case 'h':
					switch src[offset+1] {
					case 'e':
						switch src[offset+2] {
						case 'n':
							// ehrn
							if !identExtender(src[offset+3]) {
								offset += 3
								result = tokens.Identifier
								s.pseudo = tokens.When
								s.value = "when"
								break loop
							}
						case 'r':
							// where
							if src[offset+3] == 'e' && !identExtender(src[offset+4]) {
								offset += 4
								result = tokens.Identifier
								s.pseudo = tokens.Where
								s.value = "where"
								break loop
							}
						}
					case 'i':
						// while
						if src[offset+2] == 'l' && src[offset+3] == 'e' &&
							!identExtender(src[offset+4]) {
							offset += 4
							result = tokens.Identifier
							s.pseudo = tokens.While
							s.value = "while"
							break loop
						}
					}
				}
			}
			fallthrough

			// Identifier
		case 'd',
			'g', 'h', 'j',
			'k', 'm', 'n',
			'q', 's',
			'u', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E',
			'F', 'G', 'H', 'I', 'J',
			'K', 'L', 'M', 'N', 'O',
			'P', 'Q', 'R', 'S', 'T',
			'U', 'V', 'W', 'X', 'Y', 'Z',
			'_':
		identLoop:
			for {
				last := offset
				b = src[offset]
				offset++
				switch b {
				default:
					offset = last
					break identLoop
				case 'a', 'b', 'c', 'd', 'e',
					'f', 'g', 'h', 'i', 'j',
					'k', 'l', 'm', 'n', 'o',
					'p', 'q', 'r', 's', 't',
					'u', 'v', 'w', 'x', 'y', 'z',
					'A', 'B', 'C', 'D', 'E',
					'F', 'G', 'H', 'I', 'J',
					'K', 'L', 'M', 'N', 'O',
					'P', 'Q', 'R', 'S', 'T',
					'U', 'V', 'W', 'X', 'Y', 'Z',
					'_', '$',
					'0', '1', '2', '3', '4',
					'5', '6', '7', '8', '9':
					continue
				}
			}
			s.value = string(src[start:offset])
			result = tokens.Identifier
		case '0':
			if src[offset] == 'x' {
				offset++
				first := offset
				var value uint64
			hexLoop:
				for {
					last := offset
					b = src[offset]
					offset++
					switch b {
					default:
						offset = last
						break hexLoop
					case '0', '1', '2', '3', '4',
						'5', '6', '7', '8', '9':
						value = value*16 + (uint64(b) - uint64('0'))
						continue
					case 'A', 'B', 'C', 'D', 'E', 'F':
						value = value*16 + 10 + (uint64(b) - uint64('A'))
						continue
					case 'a', 'b', 'c', 'd', 'e', 'f':
						value = value*16 + 10 + (uint64(b) - uint64('a'))
						continue
					case '_':
						continue
					}
				}
				if first == offset {
					result = tokens.Invalid
					s.msg = "Invalid hex format"
				} else {
					b = src[offset]
					switch b {
					case 'u':
						offset++
						b = src[offset]
						switch b {
						default:
							result = tokens.Literal
							s.value = uint32(value)
						case 'b':
							offset++
							result = tokens.Literal
							s.value = byte(value)
						case 'l':
							offset++
							result = tokens.Literal
							s.value = uint64(value)
						}
					case 'l':
						offset++
						result = tokens.Literal
						s.value = int64(value)
					case 'i':
						offset++
						fallthrough
					default:
						result = tokens.Literal
						s.value = int32(value)
					}
				}
				switch src[offset] {
				case 'a', 'b', 'c', 'd', 'e',
					'f', 'g', 'h', 'i', 'j',
					'k', 'l', 'm', 'n', 'o',
					'p', 'q', 'r', 's', 't',
					'u', 'v', 'w', 'x', 'y',
					'z':
					result = tokens.Invalid
					s.msg = fmt.Sprintf("Extra character '%c' after literal", rune(src[offset]))
				}
				break loop
			}
			fallthrough
		case '1', '2', '3', '4',
			'5', '6', '7', '8', '9':
			value := uint64(int(b) - int('0'))
			isFloat := false
		numberLoop:
			for {
				last := offset
				b = src[offset]
				offset++
				switch b {
				default:
					offset = last
					break numberLoop
				case '0', '1', '2', '3', '4',
					'5', '6', '7', '8', '9':
					if !isFloat {
						value = value*10 + (uint64(b) - uint64('0'))
					}
				case '.':
					if src[offset] == '.' {
						offset--
						break numberLoop
					}
					isFloat = true
				}
			}
			if !isFloat {
				b = src[offset]
				switch b {
				case 'u':
					offset++
					b = src[offset]
					switch b {
					default:
						result = tokens.Literal
						s.value = uint(value)
					case 'b':
						offset++
						result = tokens.Literal
						s.value = byte(value)
					case 'l':
						offset++
						result = tokens.Literal
						s.value = value
					}
				case 'l':
					offset++
					result = tokens.Literal
					s.value = int64(value)
				case 'f':
					fvalue, err := strconv.ParseFloat(string(src[start:offset]), 32)
					offset++
					if err != nil {
						result = tokens.Invalid
						s.msg = err.Error()
					} else {
						s.value = float32(fvalue)
						result = tokens.Literal
					}
				case 'd':
					fvalue, err := strconv.ParseFloat(string(src[start:offset]), 64)
					offset++
					if err != nil {
						result = tokens.Invalid
						s.msg = err.Error()
					} else {
						s.value = fvalue
						result = tokens.Literal
					}
				case 'i':
					offset++
					fallthrough
				default:
					result = tokens.Literal
					s.value = int(value)
				}
			} else {
				var err error
				if src[offset] == 'f' {
					fvalue32, err32 := strconv.ParseFloat(string(src[start:offset]), 32)
					s.value = float32(fvalue32)
					err = err32
					result = tokens.Literal
					offset++
				} else {
					fvalue64, err64 := strconv.ParseFloat(string(src[start:offset]), 64)
					s.value = fvalue64
					err = err64
					result = tokens.Literal
					if src[offset] == 'd' {
						offset++
					}
				}
				if err != nil {
					result = tokens.Invalid
					s.msg = err.Error()
				}
			}
		case '\'':
			var value rune
			b = src[offset]
			offset++
			switch b {
			case '\\':
				value = '\\'
				b = src[offset]
				offset++
				switch b {
				case '0':
					value = '\x00'
				case 'n':
					value = '\n'
				case 'r':
					value = '\r'
				case 'b':
					value = '\b'
				case 't':
					value = '\t'
				case '\'':
					value = '\''
				default:
					result = tokens.Invalid
					s.msg = "Invalid escape"
				}
			default:
				value = rune(b)
			}
			if src[offset] != '\'' {
				result = tokens.Invalid
				s.msg = "Invalid character literal"
			} else {
				result = tokens.Literal
				offset++
			}
			s.value = value
		case '`':
			copyFrom := start + 1
			for {
				b = src[offset]
				offset++
				switch b {
				case '\n', '\r', '\\', '\x00':
					result = tokens.Invalid
					break loop
				case '`':
					s.value = string(src[copyFrom : offset-1])
					result = tokens.Identifier
					s.pseudo = tokens.Escaped
					break loop
				}
			}
		case '"':
			var value string
			copyFrom := start + 1
			for {
				b = src[offset]
				offset++
				switch b {
				case '\\':
					value += string(src[copyFrom : offset-1])
					b = src[offset]
					offset++
					switch b {
					case 'n':
						value += "\n"
						copyFrom = offset
					case 'r':
						value += "\r"
						copyFrom = offset
					case 'b':
						value += "\b"
						copyFrom = offset
					case 't':
						value += "\t"
						copyFrom = offset
					case '\\':
						value += "\\"
						copyFrom = offset
					default:
						copyFrom = offset - 1
					}
				case '"':
					value += string(src[copyFrom : offset-1])
					s.value = value
					result = tokens.Literal
					break loop
				case '\n', '\r', 0:
					result = tokens.Invalid
					s.msg = "Unterminated string"
					break loop
				}
			}
		}
		break loop
	}
	s.offset = offset
	s.start = start
	s.end = offset
	s.line = line
	return result
}
