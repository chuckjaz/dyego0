package parser_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "go/token"

    "dyego0/ast"
    "dyego0/parser"
)

var _ = Describe("parser", func() {
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
        b := func() ast.Builder {
            b := ast.NewBuilder(scanner(""))
            b.PushContext()
            return b
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
})

func scanner(text string) *parser.Scanner {
    return parser.NewScanner(append([]byte(text), 0), 0)
}

func parse(text string) ast.Element {
    p := parser.NewParser(scanner(text))
    r := p.Parse()
    Expect(p.Errors()).To(BeNil())
    return r
}

func TestErrors(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Errors Suite")
}
