package parser

import "strconv"

import "go/token"
import "dyego0/tokens"

const (
	// InternalScan enables internal identifiers
	InternalScan = 1 << iota
)

// Scanner is a Dyego scanner
type Scanner struct {
	src    []byte
	offset int
	line   int
	start  int
	end    int
	msg    string
	flags  int
	value  interface{}
}

// NewScanner creates a scanner
func NewScanner(src []byte, flags int) *Scanner {
	length := len(src)
	if length == 0 || src[length-1] != 0 {
		panic("NewScanner: src must be null terminated")
	}
	return &Scanner{src: src, line: 1, flags: flags}
}

// Clone preserves a copy of the scanner at the current state which can then be used
// for backtracking, if necessary by using the returned instance instead of the
// instance that was moved forward.
func (s *Scanner) Clone() *Scanner {
	return &Scanner{src: s.src, offset: s.offset, line: s.line, start: s.start,
		end: s.end, msg: s.msg, flags: s.flags, value: s.value}
}

// Line is the current line of the scanner
func (s *Scanner) Line() int {
	return s.line
}

// Start is the start of the current token
func (s *Scanner) Start() token.Pos {
	return token.Pos(s.start)
}

// End is the end of the current token
func (s *Scanner) End() token.Pos {
	return token.Pos(s.end)
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

// IsInternal reports if the scanning an internal file
func (s *Scanner) IsInternal() bool {
	return s.flags&InternalScan != 0
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

// Next moves the scanner to the next token
func (s *Scanner) Next() tokens.Token {
	s.end = s.offset
	offset := s.offset
	start := s.offset
	src := s.src
	line := s.line
	result := tokens.Invalid
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
			continue

		case '+':
			result = tokens.Add
		case '-':
			if src[offset] == '>' {
				offset++
				result = tokens.Arrow
			} else {
				result = tokens.Sub
			}
		case '*':
			result = tokens.Mult
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
				continue
			}
			result = tokens.Div
		case '%':
			result = tokens.Rem
		case '(':
			result = tokens.LParen
		case ')':
			result = tokens.RParen
		case '[':
			result = tokens.LBrack
		case ']':
			result = tokens.RBrack
		case '{':
			result = tokens.LBrace
		case '}':
			result = tokens.RBrace
		case '!':
			if src[offset] == '=' {
				offset++
				result = tokens.NotEqual
			} else if src[offset] == '!' && s.IsInternal() {
				offset++
				result = tokens.Platform
			} else {
				result = tokens.Not
			}
		case ',':
			result = tokens.Comma
		case '=':
			result = tokens.Equal
		case ':':
			if src[offset] == '=' {
				offset++
				result = tokens.Assign
			} else {
				result = tokens.Colon
			}
		case '?':
			result = tokens.Question
		case ';':
			result = tokens.Semi
		case '.':
			if src[offset] == '.' {
				offset++
				result = tokens.Range
			} else {
				result = tokens.Dot
			}

		case '&':
			if src[offset] == '&' {
				offset++
				result = tokens.LogicalAnd
			} else {
				result = tokens.Invalid
			}

		case '|':
			if src[offset] == '|' {
				offset++
				result = tokens.LogicalOr
			} else {
				result = tokens.Invalid
			}

		case '>':
			if src[offset] == '=' {
				offset++
				result = tokens.GreaterThanEqual
			} else {
				result = tokens.GreaterThan
			}

		case '<':
			if src[offset] == '=' {
				offset++
				result = tokens.LessThanEqual
			} else {
				result = tokens.LessThan
			}

		case '$':
			if !s.IsInternal() {
				result = tokens.Invalid
				s.msg = "invalid identifier"
				break loop
			}
			fallthrough

		// Reserved words
		case 'a', 'b', 'c', 'd', 'e', 'f', 'i', 'l', 'm', 'o', 'p', 'r', 's', 't', 'v':
			switch b {
			case 'a':
				if src[offset] == 's' && !identExtender(src[offset+1]) {
					offset++
					result = tokens.As
					break loop
				}
			case 'b':
				switch src[offset] {
				case 'o':
					// boolean
					if src[offset+1] == 'o' && src[offset+2] == 'l' &&
						src[offset+3] == 'e' && src[offset+4] == 'a' &&
						src[offset+5] == 'n' && !identExtender(src[offset+6]) {
						offset += 6
						result = tokens.Boolean
						break loop
					}
				case 'r':
					// break
					if src[offset+1] == 'e' && src[offset+2] == 'a' &&
						src[offset+3] == 'k' && !identExtender(src[offset+4]) {
						offset += 4
						result = tokens.Break
						break loop
					}
				}
			case 'c':
				switch src[offset] {
				case 'a':
					// case
					if src[offset+1] == 's' && src[offset+2] == 'e' && !identExtender(src[offset+3]) {
						offset += 3
						result = tokens.Case
						break loop
					}
				case 'o':
					if src[offset+1] == 'n' {
						switch src[offset+2] {
						case 's':
							// constraint
							if src[offset+3] == 't' && src[offset+4] == 'r' && src[offset+5] == 'a' &&
								src[offset+6] == 'i' && src[offset+7] == 'n' && src[offset+8] == 't' &&
								!identExtender(src[offset+9]) {
								offset += 9
								result = tokens.Constraint
								break loop
							}
						case 't':
							// continue
							if src[offset+3] == 'i' && src[offset+4] == 'n' && src[offset+5] == 'u' &&
								src[offset+6] == 'e' && !identExtender(src[offset+7]) {
								offset += 7
								result = tokens.Continue
								break loop
							}
						}
					}
				}
			case 'd':
				// default
				if src[offset] == 'e' && src[offset+1] == 'f' && src[offset+2] == 'a' &&
					src[offset+3] == 'u' && src[offset+4] == 'l' && src[offset+5] == 't' &&
					!identExtender(src[offset+6]) {
					offset += 6
					result = tokens.Default
					break loop
				}
			case 'e':
				switch src[offset] {
				case 'l':
					// else
					if src[offset+1] == 's' && src[offset+2] == 'e' && !identExtender(src[offset+3]) {
						offset += 3
						result = tokens.Else
						break loop
					}
				case 'n':
					// enum
					if src[offset+1] == 'u' && src[offset+2] == 'm' && !identExtender(src[offset+3]) {
						offset += 3
						result = tokens.Enum
						break loop
					}
				}
			case 'f':
				switch src[offset] {
				case 'a':
					// false
					if src[offset+1] == 'l' && src[offset+2] == 's' &&
						src[offset+3] == 'e' && !identExtender(src[offset+4]) {
						offset += 4
						result = tokens.False
						break loop
					}
				case 'l':
					// float
					if src[offset+1] == 'o' && src[offset+2] == 'a' &&
						src[offset+3] == 't' && !identExtender(src[offset+4]) {
						offset += 4
						result = tokens.Float
						break loop
					}
				case 'o':
					// for
					if src[offset+1] == 'r' && !identExtender(src[offset+2]) {
						offset += 2
						result = tokens.For
						break loop
					}
				}
			case 'i':
				switch src[offset] {
				case 'f':
					// if
					if !identExtender(src[offset+1]) {
						offset++
						result = tokens.If
						break loop
					}
				case 'n':
					if src[offset+1] == 't' {
						// interface
						if src[offset+2] == 'e' && src[offset+3] == 'r' && src[offset+4] == 'f' && src[offset+5] == 'a' &&
							src[offset+6] == 'c' && src[offset+7] == 'e' && !identExtender(src[offset+8]) {
							offset += 8
							result = tokens.Interface
							break loop
						}

						// int
						if !identExtender(src[offset+2]) {
							offset += 2
							result = tokens.Int
							break loop
						}
					}
					// in
					if !identExtender(src[offset+1]) {
						offset++
						result = tokens.In
						break loop
					}
				}
			case 'l':
				// let
				if src[offset] == 'e' && src[offset+1] == 't' && !identExtender(src[offset+2]) {
					offset += 2
					result = tokens.Let
					break loop
				}
			case 'm':
				switch src[offset] {
				case 'a':
					// match
					if src[offset+1] == 't' && src[offset+2] == 'c' &&
						src[offset+3] == 'h' && !identExtender(src[offset+4]) {
						offset += 4
						result = tokens.Match
						break loop
					}
				case 'e':
					// method
					if src[offset+1] == 't' && src[offset+2] == 'h' &&
						src[offset+3] == 'o' && src[offset+4] == 'd' && !identExtender(src[offset+5]) {
						offset += 5
						result = tokens.Method
						break loop
					}
				}
			case 'o':
				// operator
				if src[offset] == 'p' && src[offset+1] == 'e' && src[offset+2] == 'r' &&
					src[offset+3] == 'a' && src[offset+4] == 't' && src[offset+5] == 'o' &&
					src[offset+6] == 'r' && !identExtender(src[offset+7]) {
					offset += 7
					result = tokens.Operator
					break loop
				}
			case 'p':
				// property
				if src[offset] == 'r' && src[offset+1] == 'o' && src[offset+2] == 'p' &&
					src[offset+3] == 'e' && src[offset+4] == 'r' && src[offset+5] == 't' &&
					src[offset+6] == 'y' && !identExtender(src[offset+7]) {
					offset += 7
					result = tokens.Property
					break loop
				}
			case 'r':
				// return
				if src[offset] == 'e' {
					switch src[offset+1] {
					case 'c':
						// record
						if src[offset+2] == 'o' && src[offset+3] == 'r' && src[offset+4] == 'd' &&
							!identExtender(src[offset+5]) {
							offset += 5
							result = tokens.Record
							break loop
						}
					case 't':
						if src[offset+2] == 'u' && src[offset+3] == 'r' && src[offset+4] == 'n' &&
							!identExtender(src[offset+5]) {
							offset += 5
							result = tokens.Return
							break loop
						}
					}
				}
			case 's':
				switch src[offset] {
				case 't':
					// string
					if src[offset+1] == 'r' && src[offset+2] == 'i' &&
						src[offset+3] == 'n' && src[offset+4] == 'g' && !identExtender(src[offset+5]) {
						offset += 5
						result = tokens.String
						break loop
					}
				case 'w':
					// switch
					if src[offset+1] == 'i' && src[offset+2] == 't' && src[offset+3] == 'c' &&
						src[offset+4] == 'h' && !identExtender(src[offset+5]) {
						offset += 5
						result = tokens.Switch
						break loop
					}
				}
			case 't':
				switch src[offset] {
				case 'r':
					// true
					if src[offset+1] == 'u' && src[offset+2] == 'e' && !identExtender(src[offset+3]) {
						offset += 3
						result = tokens.True
						break loop
					}
				case 'y':
					// type
					if src[offset+1] == 'p' && src[offset+2] == 'e' && !identExtender(src[offset+3]) {
						offset += 3
						result = tokens.Type
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
						// value
						if src[offset+2] == 'u' && src[offset+3] == 'e' && !identExtender(src[offset+4]) {
							offset += 4
							result = tokens.Value
							break loop
						}
					}
				case 'o':
					// void
					if src[offset+1] == 'i' && src[offset+2] == 'd' && !identExtender(src[offset+3]) {
						offset += 3
						result = tokens.Void
						break loop
					}
				}
			}
			fallthrough

			// Identifier
		case 'g', 'h', 'j',
			'k', 'n',
			'q',
			'u', 'w', 'x', 'y', 'z',
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
		case '0', '1', '2', '3', '4',
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
				case 'b':
					offset++
					result = tokens.LiteralByte
					s.value = byte(value)
				case 'u':
					offset++
					b = src[offset]
					switch b {
					default:
						result = tokens.LiteralUInt
						s.value = uint32(value)
					case 'l':
						offset++
						result = tokens.LiteralULong
						s.value = value
					}
				case 'l':
					offset++
					result = tokens.LiteralLong
					s.value = int64(value)
				case 'i':
					offset++
					fallthrough
				case 'f':
					fvalue, err := strconv.ParseFloat(string(src[start:offset]), 32)
					offset++
					if err != nil {
						result = tokens.Invalid
						s.msg = err.Error()
					} else {
						s.value = float32(fvalue)
						result = tokens.LiteralFloat
					}
				case 'd':
					fvalue, err := strconv.ParseFloat(string(src[start:offset]), 64)
					offset++
					if err != nil {
						result = tokens.Invalid
						s.msg = err.Error()
					} else {
						s.value = fvalue
						result = tokens.LiteralFloat
					}
				default:
					result = tokens.LiteralInt
					s.value = int(value)
				}
			} else {
				var err error
				if src[offset] == 'f' {
					fvalue32, err32 := strconv.ParseFloat(string(src[start:offset]), 32)
					s.value = float32(fvalue32)
					err = err32
					result = tokens.LiteralFloat
					offset++
				} else {
					fvalue64, err64 := strconv.ParseFloat(string(src[start:offset]), 64)
					s.value = fvalue64
					err = err64
					result = tokens.LiteralDouble
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
				result = tokens.LiteralRune
				offset++
			}
			s.value = value
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
					result = tokens.LiteralString
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
	s.line = line
	return result
}
