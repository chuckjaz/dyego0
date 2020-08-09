package parser

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"go/token"

	"dyego0/assert"
	"dyego0/ast"
	"dyego0/scanner"
)

var _ = Describe("parser", func() {
	b := func() ast.Builder {
		b := ast.NewBuilder(scan(""))
		b.PushContext()
		return b
	}
	Describe("literals", func() {
		It("can parse a rune", func() {
			l, ok := parse("'a'").(ast.LiteralRune)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal('a'))
		})
		It("can parse a byte", func() {
			l, ok := parse("42b").(ast.LiteralByte)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(byte(42)))
		})
		It("can parse an int", func() {
			l, ok := parse("123").(ast.LiteralInt)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(123))
		})
		It("can parse a uint", func() {
			l, ok := parse("42u").(ast.LiteralUInt)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(uint(42)))
		})
		It("can parse a long", func() {
			l, ok := parse("42l").(ast.LiteralLong)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(int64(42)))
		})
		It("can parse a unsigned long", func() {
			l, ok := parse("42ul").(ast.LiteralULong)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(uint64(42)))
		})
		It("can parse a float", func() {
			l, ok := parse("1.0f").(ast.LiteralFloat)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(BeNumerically("~", float32(1.0)))
		})
		It("can parse a double", func() {
			l, ok := parse("1.0").(ast.LiteralDouble)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(1.0))
		})
		It("can parse a string", func() {
			l, ok := parse("\"a\"").(ast.LiteralString)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal("a"))
		})
		It("can parse true", func() {
			l, ok := parse("true").(ast.LiteralBoolean)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(true))
		})
		It("can parse false", func() {
			l, ok := parse("false").(ast.LiteralBoolean)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(false))
		})
		It("can parse a parenthesised expression", func() {
			l, ok := parse("(10)").(ast.LiteralInt)
			Expect(ok).To(Equal(true))
			Expect(l.Value()).To(Equal(10))
		})
	})
	Describe("simple expression", func() {
		It("can parse a selection", func() {
			l, ok := parse("a.b").(ast.Selection)
			Expect(ok).To(Equal(true))
			t, ok := l.Target().(ast.Name)
			Expect(ok).To(Equal(true))
			Expect(t.Text()).To(Equal("a"))
			Expect(l.Member().Text()).To(Equal("b"))
			Expect(l.Start()).To(Equal(token.Pos(0)))
			Expect(l.End()).To(Equal(token.Pos(3)))
		})
		It("can parse a call", func() {
			c, ok := parse("a(10, b = 2)").(ast.Call)
			Expect(ok).To(Equal(true))
			n, ok := c.Target().(ast.Name)
			Expect(ok).To(Equal(true))
			Expect(n.Text()).To(Equal("a"))
			arguments := c.Arguments()
			Expect(len(arguments)).To(Equal(2))
			num, ok := arguments[0].(ast.LiteralInt)
			Expect(ok).To(Equal(true))
			Expect(num.Value()).To(Equal(10))
			na, ok := arguments[1].(ast.NamedArgument)
			Expect(ok).To(Equal(true))
			Expect(na.Name().Text()).To(Equal("b"))
			num2, ok := na.Value().(ast.LiteralInt)
			Expect(ok).To(Equal(true))
			Expect(num2.Value()).To(Equal(2))
		})
	})
	Describe("lambda", func() {
		lambda := func(source string) ast.Lambda {
			l, ok := parse(source).(ast.Lambda)
			Expect(ok).To(Equal(true))
			return l
		}
		expectNumber := func(element ast.Element, value int) {
			n, ok := element.(ast.LiteralInt)
			Expect(ok).To(Equal(true))
			Expect(n.Value()).To(Equal(value))
		}
		expectName := func(element ast.Element, value string) {
			l, ok := element.(ast.Name)
			Expect(ok).To(Equal(true))
			Expect(l.Text()).To(Equal(value))
		}
		expectNil := func(element interface{}) {
			Expect(element).To(BeNil())
		}
		p := func(name string) ast.Parameter {
			return b().Parameter(b().Name(name), nil, nil)
		}
		pd := func(name string) ast.Parameter {
			return b().Parameter(b().Name(name), nil, b().LiteralInt(42))
		}
		pt := func(name string, typ string) ast.Parameter {
			return b().Parameter(b().Name(name), b().Name(typ), nil)
		}
		ptd := func(name string, typ string) ast.Parameter {
			return b().Parameter(b().Name(name), b().Name(typ), b().LiteralInt(42))
		}
		tp := func(name string) ast.TypeParameter {
			return b().TypeParameter(b().Name(name), nil)
		}
		tpc := func(name, constraint string) ast.TypeParameter {
			return b().TypeParameter(b().Name(name), b().Name(constraint))
		}
		w := func(left, right string) ast.Where {
			return b().Where(b().Name(left), b().Name(right))
		}
		expectParameter := func(parameter, expected ast.Parameter) {
			Expect(parameter.Name().Text()).To(Equal(expected.Name().Text()))
			if expected.Type() != nil {
				n := expected.Type().(ast.Name)
				expectName(parameter.Type(), n.Text())
			}
			if expected.Default() != nil {
				expectNumber(parameter.Default(), 42)
			}
		}
		expectParameters := func(parameters []ast.Parameter, expected ...ast.Parameter) {
			Expect(len(parameters)).To(Equal(len(expected)))
			for i := 0; i < len(parameters); i++ {
				expectParameter(parameters[i], expected[i])
			}
		}
		expectTypeParameter := func(parameter, expected ast.TypeParameter) {
			Expect(parameter.Name().Text()).To(Equal(expected.Name().Text()))
			if expected.Constraint() != nil {
				c := expected.Constraint().(ast.Name)
				expectName(parameter.Constraint(), c.Text())
			}
		}
		expectTypeParameters := func(parameters ast.TypeParameters, expected ...ast.TypeParameter) {
			Expect(len(parameters.Parameters())).To(Equal(len(expected)))
			for i := 0; i < len(expected); i++ {
				expectTypeParameter(parameters.Parameters()[i], expected[i])
			}
		}
		expectWhere := func(where, expected ast.Where) {
			left := expected.Left().(ast.Name).Text()
			right := expected.Right().(ast.Name).Text()
			expectName(where.Left(), left)
			expectName(where.Right(), right)
		}
		expectWheres := func(parameters ast.TypeParameters, expected ...ast.Where) {
			Expect(len(parameters.Wheres())).To(Equal(len(expected)))
			for i := 0; i < len(expected); i++ {
				expectWhere(parameters.Wheres()[i], expected[i])
			}
		}
		It("can parse an empty lambda", func() {
			l := lambda("{}")
			expectNil(l.TypeParameters())
			expectNil(l.Parameters())
			expectNil(l.Body())
		})
		It("can parse a simple lambda expression", func() {
			l := lambda("{ 42 }")
			expectNumber(l.Body(), 42)
		})
		It("can parse lambda with a parameter", func() {
			l := lambda("{ a -> a }")
			expectName(l.Body(), "a")
			expectParameters(l.Parameters(), p("a"))
		})
		It("can parse lambda with multiple parameters", func() {
			l := lambda("{ a, b -> 42 }")
			expectNumber(l.Body(), 42)
			expectParameters(l.Parameters(), p("a"), p("b"))
		})
		It("can parse lambda with a default parameter value", func() {
			l := lambda("{ a = 42 -> a }")
			expectName(l.Body(), "a")
			expectParameters(l.Parameters(), pd("a"))
		})
		It("can parse a lambda with a typed defualt parameter", func() {
			l := lambda("{ a: Int = 42 -> a }")
			expectName(l.Body(), "a")
			expectParameters(l.Parameters(), ptd("a", "Int"))
		})
		It("can parse a lambda with a type parameter", func() {
			l := lambda("{ A | a: A -> a }")
			expectName(l.Body(), "a")
			expectTypeParameters(l.TypeParameters(), tp("A"))
			expectParameters(l.Parameters(), pt("a", "A"))
		})
		It("can parse a lambda with multimple type parameters", func() {
			l := lambda("{ A, B | a: A, b: B -> 42 }")
			expectNumber(l.Body(), 42)
			expectTypeParameters(l.TypeParameters(), tp("A"), tp("B"))
			expectParameters(l.Parameters(), pt("a", "A"), pt("b", "B"))
		})
		It("can parse a type parameter with a constraint", func() {
			l := lambda("{ A: Int | a: A -> a }")
			expectName(l.Body(), "a")
			expectTypeParameters(l.TypeParameters(), tpc("A", "Int"))
		})
		It("can parse a type where clause", func() {
			l := lambda("{ A where A = Int | a: A -> a }")
			expectName(l.Body(), "a")
			expectWheres(l.TypeParameters(), w("A", "Int"))
		})
		It("can parser multiple where classes", func() {
			l := lambda("{ A, B where A = Int where B = A | a: A -> a }")
			expectName(l.Body(), "a")
			expectWheres(l.TypeParameters(), w("A", "Int"), w("B", "A"))
		})
	})
	Describe("vocabulary", func() {
		vocab := func(source string) ast.VocabularyLiteral {
			l, ok := parse("let v = " + source).(ast.LetDefinition)
			Expect(ok).To(Equal(true))
			v, ok := l.Value().(ast.VocabularyLiteral)
			return v
		}
		It("can parse an empty vocabulary", func() {
			v := vocab("<| |>")
			Expect(v.Members()).To(BeNil())
		})
		Describe("embedding", func() {
			ve := func(source string) ast.VocabularyEmbedding {
				v := vocab(source)
				m := v.Members()
				Expect(len(m)).To(Equal(1))
				return m[0].(ast.VocabularyEmbedding)
			}
			expectNames := func(name []ast.Name, expected ...string) {
				Expect(len(name)).To(Equal(len(expected)))
				for i := range name {
					Expect(name[i].Text()).To(Equal(expected[i]))
				}
			}
			It("can parse a vocabulary embedding", func() {
				e := ve("<| ...Other |>")
				expectNames(e.Name(), "Other")
			})
			It("can parse a bocabulary embedding reference", func() {
				e := ve("<| ...a::b |>")
				expectNames(e.Name(), "a", "b")
			})
		})
		Describe("operator", func() {
			op := func(source string) ast.VocabularyOperatorDeclaration {
				v := vocab("<| " + source + " |>")
				Expect(len(v.Members())).To(Equal(1))
				o, ok := v.Members()[0].(ast.VocabularyOperatorDeclaration)
				Expect(ok).To(Equal(true))
				return o
			}
			It("can parse an infix operator", func() {
				o := op("infix operator `+` left")
				Expect(o.Placement()).To(Equal(ast.Infix))
				Expect(len(o.Names())).To(Equal(1))
				Expect(o.Names()[0].Text()).To(Equal("+"))
				Expect(o.Associativity()).To(Equal(ast.Left))
			})
			It("can parse an prefix operator", func() {
				o := op("prefix operator `+` left")
				Expect(o.Placement()).To(Equal(ast.Prefix))
				Expect(len(o.Names())).To(Equal(1))
				Expect(o.Names()[0].Text()).To(Equal("+"))
				Expect(o.Associativity()).To(Equal(ast.Left))
			})
			It("can parse an postfix operator", func() {
				o := op("postfix operator `+` left")
				Expect(o.Placement()).To(Equal(ast.Postfix))
				Expect(len(o.Names())).To(Equal(1))
				Expect(o.Names()[0].Text()).To(Equal("+"))
				Expect(o.Associativity()).To(Equal(ast.Left))
			})
			It("can parse right associative operator", func() {
				o := op("infix operator `+` right")
				Expect(o.Associativity()).To(Equal(ast.Right))
			})
			It("can parse multiple names in operator declaration", func() {
				o := op("infix operator (`+`, `-`) right")
				Expect(len(o.Names())).To(Equal(2))
				Expect(o.Names()[0].Text()).To(Equal("+"))
				Expect(o.Names()[1].Text()).To(Equal("-"))
			})
			It("can parse mutlpile operators", func() {
				v := vocab("<| infix operator a left, infix operator b left |>")
				Expect(len(v.Members())).To(Equal(2))
			})
		})
	})
	e := func(source string) ast.Element {
		actualSource := "...dyego," + source
		element := parse(actualSource)
		sequence := element.(ast.Sequence)
		return sequence.Right()
	}
	Describe("prefix expressions", func() {
		expectName := func(n ast.Element, name string) {
			nameElement, ok := n.(ast.Name)
			Expect(ok).To(BeTrue())
			Expect(nameElement.Text()).To(Equal(name))
		}
		i := func(n ast.Element, value int) {
			v, ok := n.(ast.LiteralInt)
			Expect(ok).To(BeTrue())
			Expect(v.Value()).To(Equal(value))
		}
		expectOp := func(e ast.Element) (ast.Element, ast.Element, []ast.Element) {
			call, ok := e.(ast.Call)
			Expect(ok).To(BeTrue())
			selection, ok := call.Target().(ast.Selection)
			Expect(ok).To(BeTrue())
			return selection.Member(), selection.Target(), call.Arguments()
		}
		expectUnaryOp := func(e ast.Element, name string) ast.Element {
			member, target, arguments := expectOp(e)
			expectName(member, name)
			Expect(len(arguments)).To(Equal(0))
			return target
		}
		expectBinaryOp := func(e ast.Element, name string) (ast.Element, ast.Element) {
			member, target, arguments := expectOp(e)
			expectName(member, name)
			Expect(len(arguments)).To(Equal(1))
			return target, arguments[0]
		}
		It("can parse a prefix expression", func() {
			v := e("+1")
			t := expectUnaryOp(v, "+")
			i(t, 1)
		})
		It("can parse a binary expression", func() {
			v := e("1 + 2")
			l, r := expectBinaryOp(v, "+")
			i(l, 1)
			i(r, 2)
		})
		It("can distinquish precedence", func() {
			v := e("1 * 2 + 3 * 4")
			lm, rm := expectBinaryOp(v, "+")
			e1, e2 := expectBinaryOp(lm, "*")
			e3, e4 := expectBinaryOp(rm, "*")
			i(e1, 1)
			i(e2, 2)
			i(e3, 3)
			i(e4, 4)
		})
		It("can handle mix of operator types", func() {
			v := e("++1++ + ++2++")
			l, r := expectBinaryOp(v, "+")
			expectUnaryOp(expectUnaryOp(l, "++"), "postfix ++")
			expectUnaryOp(expectUnaryOp(r, "++"), "postfix ++")
		})
		It("can handle multiple prefix/postfix operators", func() {
			v := e("++ ++ 1 ++ ++")
			expectUnaryOp(expectUnaryOp(expectUnaryOp(expectUnaryOp(v, "++"), "++"), "postfix ++"), "postfix ++")
		})
	})
})

func scan(text string) *scanner.Scanner {
	return scanner.NewScanner(append([]byte(text), 0), 0)
}

var dyego0VocabularySource = strings.ReplaceAll(`
let dyego = <| 
  postfix operator (@++@, @--@, @?.@, @?@) right,
  prefix operator (@+@, @-@, @--@, @++@) right,
  infix operator (@as@, @as?@) left,
  infix operator (@*@, @/@, @%@) left,
  infix operator (@+@, @-@) left,
  infix operator @..@ left,
  infix operator @?:@ left,
  infix operator (@in@, @!in@, @is@, @!is@) left,
  infix operator (@<@, @>@, @>=@, @<=@) left,
  infix operator (@==@, @!=@) left,
  infix operator @&&@ left,
  infix operator @||@ left,
  infix operator (@=@, @+=@, @*=@, @/=@, @%=@) right
|>`, "@", "`")

func lineMapOf(source string) []int {
	result := []int{0}
	for i, ch := range source {
		if ch == '\n' {
			result = append(result, i+1)
		}
	}
	return append(result, len(source))
}

func lineColumnOf(lineMap []int, offset int) (int, int) {
	line0 := sort.SearchInts(lineMap, offset)
	if lineMap[line0] != offset {
		line0--
	}
	if line0 < len(lineMap) {
		return line0 + 1, offset - lineMap[line0] + 1
	} else if lineMap != nil {
		return line0, offset - lineMap[line0-1] + 1
	} else {
		return 1, offset + 1
	}
}

func lineLenOf(lineMap []int, line int) int {
	if lineMap == nil {
		return 0
	}
	return lineMap[line] - lineMap[line-1]
}

func printErrors(errors []ast.Error, source string) {
	if errors != nil {
		lineMap := lineMapOf(source)
		for _, err := range errors {
			offset := int(err.Start())
			endOffset := int(err.End())
			line, col := lineColumnOf(lineMap, offset)
			errorLen := endOffset - offset
			lineLength := lineLenOf(lineMap, line)
			if lineLength-col < errorLen {
				errorLen = lineLength - col
			}
			println(err.Message())
			println(source[lineMap[line-1]:lineMap[line]])
			println(fmt.Sprintf("%s%s", strings.Repeat(" ", col-1), strings.Repeat("^", errorLen)))
		}
	}
}

func parseVocabulary(source string) vocabulary {
	p := NewParser(scan(source), newVocabularyScope())
	element := p.Parse()
	printErrors(p.Errors(), source)
	assert.Assert(len(p.Errors()) == 0, "dyego0 vocabulary source has errors %#v", p.Errors())
	let := element.(ast.LetDefinition)
	vl := let.Value().(ast.VocabularyLiteral)
	emptyScope := newVocabularyScope()
	v, errors := buildVocabulary(emptyScope, vl)
	assert.Assert(len(errors) == 0, "dyego0 vocabulary source has errors %#v", errors)
	return v
}

var defaultScope vocabularyScope

func parse(text string) ast.Element {
	p := NewParser(scan(text), defaultScope)
	r := p.Parse()
	printErrors(p.Errors(), text)
	Expect(p.Errors()).To(BeNil())
	return r
}

func TestErrors(t *testing.T) {
	dyegoVocabulary := parseVocabulary(dyego0VocabularySource)
	scope := newVocabularyScope()
	scope.members["dyego"] = dyegoVocabulary
	defaultScope = scope

	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}
