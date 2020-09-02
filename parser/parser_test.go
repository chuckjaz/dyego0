package parser

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

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
	i := func(n ast.Element, value int) {
		v, ok := n.(ast.LiteralInt)
		Expect(ok).To(BeTrue())
		Expect(v.Value()).To(Equal(value))
	}
	n := func(n ast.Element, value string) {
		v, ok := n.(ast.Name)
		Expect(ok).To(BeTrue())
		Expect(v.Text()).To(Equal(value))
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
	named := func(name string, o ast.ObjectInitializer) ast.Element {
		for _, member := range o.Members() {
			namedMember, ok := member.(ast.NamedMemberInitializer)
			if ok {
				if namedMember.Name().Text() == name {
					return namedMember.Value()
				}
			}
		}
		assert.Fail("No member named %s found", name)
		return nil
	}
	Describe("object initializer", func() {
		obj := func(text string) ast.ObjectInitializer {
			r := parse(text)
			o, ok := r.(ast.ObjectInitializer)
			Expect(ok).To(BeTrue())
			return o
		}
		ro := func(text string) ast.ObjectInitializer {
			obj := obj(text)
			Expect(obj.Mutable()).To(BeFalse())
			return obj
		}
		mo := func(text string) ast.ObjectInitializer {
			obj := obj(text)
			Expect(obj.Mutable()).To(BeTrue())
			return obj
		}
		Describe("read only", func() {
			It("can parse one field", func() {
				o := ro("[a: 1]")
				v := named("a", o)
				i(v, 1)
			})
			It("can parse two fields", func() {
				o := ro("[a: 1, b: 2]")
				i(named("a", o), 1)
				i(named("b", o), 2)
			})
			It("can parse typed", func() {
				o := ro("[<A> a: 1]")
				i(named("a", o), 1)
				n(o.Type(), "A")
			})
			It("can parse simpilified member", func() {
				o := ro("[:a]")
				n(named("a", o), "a")
			})
			It("can parse a member spread", func() {
				o := ro("[...a]")
				s := o.Members()[0].(ast.SpreadMemberInitializer)
				n(s.Target(), "a")
			})
		})
		Describe("mutable", func() {
			It("can parse one field", func() {
				o := mo("[! a: 1 !]")
				v := named("a", o)
				i(v, 1)
			})
			It("can parse two fields", func() {
				o := mo("[! a: 1, b: 2 !]")
				i(named("a", o), 1)
				i(named("b", o), 2)
			})
			It("can parse typed", func() {
				o := mo("[! <A> a: 1 !]")
				i(named("a", o), 1)
				n(o.Type(), "A")
			})
		})
	})
	Describe("object initializer", func() {
		arr := func(text string) ast.ArrayInitializer {
			r := parse(text)
			o, ok := r.(ast.ArrayInitializer)
			Expect(ok).To(BeTrue())
			return o
		}
		ra := func(text string) ast.ArrayInitializer {
			a := arr(text)
			Expect(a.Mutable()).To(BeFalse())
			return a
		}
		ma := func(text string) ast.ArrayInitializer {
			a := arr(text)
			Expect(a.Mutable()).To(BeTrue())
			return a
		}
		obj := func(e ast.Element) ast.ObjectInitializer {
			o, ok := e.(ast.ObjectInitializer)
			Expect(ok).To(BeTrue())
			return o
		}
		Describe("read only", func() {
			It("can parse an single element array", func() {
				a := ra("[ 1 ]")
				i(a.Elements()[0], 1)
			})
			It("can parse a multi-element array", func() {
				a := ra("[ 1, 2, 3]")
				i(a.Elements()[0], 1)
				i(a.Elements()[1], 2)
				i(a.Elements()[2], 3)
			})
			It("can parse a typed array", func() {
				a := ra("[<A> 1]")
				i(a.Elements()[0], 1)
				n(a.Type(), "A")
			})
			It("can parse a nested array", func() {
				a := ra("[ [ a: 1] ]")
				i(named("a", obj(a.Elements()[0])), 1)
			})
			It("can parse mixed array", func() {
				parse(`
                  ...dyego

                  [<Sphere[]>
                    [ center: [ x: -1.0, y: 1.0 - t/10.0, z: 3.0 ]
                      radius: 0.3
                      color: red ]
                    [ center: [ x: 0.0, y: 1.0 - t/10.0, z: 3.0 - t/4.0 ]
                      radius: 0.3
                      color: red ]
                  ]`)
			})
		})
		Describe("mutable", func() {
			It("can parse an single element array", func() {
				a := ma("[! 1 !]")
				i(a.Elements()[0], 1)
			})
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
			na := func(e ast.Element) ast.NamedArgument {
				r, ok := e.(ast.NamedArgument)
				Expect(ok).To(BeTrue())
				return r
			}
			c, ok := parse("a(10, b: 2, :c)").(ast.Call)
			Expect(ok).To(Equal(true))
			n(c.Target(), "a")
			arguments := c.Arguments()
			Expect(len(arguments)).To(Equal(3))
			i(arguments[0], 10)
			nb := na(arguments[1])
			n(nb.Name(), "b")
			i(nb.Value(), 2)
			nc := na(arguments[2])
			n(nc.Name(), "c")
			n(nc.Value(), "c")
		})
		It("can parse an index", func() {
			c, ok := parse("a[1]").(ast.Call)
			Expect(ok).To(BeTrue())
			s, ok := c.Target().(ast.Selection)
			Expect(ok).To(BeTrue())
			n(s.Target(), "a")
			n(s.Member(), "get")
			i(c.Arguments()[0], 1)
		})
		It("can parse an index assignment", func() {
			c, ok := parse("a[1] = 2").(ast.Call)
			Expect(ok).To(BeTrue())
			s, ok := c.Target().(ast.Selection)
			Expect(ok).To(BeTrue())
			n(s.Target(), "a")
			n(s.Member(), "set")
			i(c.Arguments()[0], 1)
			i(c.Arguments()[1], 2)
		})
	})
	Describe("lambda", func() {
		lambda := func(source string) ast.Lambda {
			l, ok := parse(source).(ast.Lambda)
			Expect(ok).To(Equal(true))
			return l
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
		It("can parse type parameter with a trailing comma", func() {
			l := lambda("{ A, | a: A -> a }")
			expectTypeParameters(l.TypeParameters(), tp("A"))
		})
	})
	Describe("statements", func() {
		It("can parser a loop", func() {
			_, ok := parse("loop { 1 }").(ast.Loop)
			Expect(ok).To(BeTrue())
		})
		It("can parse a labeled loop", func() {
			l, ok := parse("loop loop { 1 }").(ast.Loop)
			Expect(ok).To(BeTrue())
			expectName(l.Label(), "loop")
		})
		It("can parse break", func() {
			_, ok := parse("break").(ast.Break)
			Expect(ok).To(BeTrue())
		})
		It("can parse labeled break", func() {
			b, ok := parse("break loop").(ast.Break)
			Expect(ok).To(BeTrue())
			expectName(b.Label(), "loop")
		})
		It("can parse continue", func() {
			_, ok := parse("continue").(ast.Continue)
			Expect(ok).To(BeTrue())
		})
		It("can parse labeled continue", func() {
			c, ok := parse("continue loop").(ast.Continue)
			Expect(ok).To(BeTrue())
			expectName(c.Label(), "loop")
		})
		It("can parse a return statement", func() {
			_, ok := parse("return").(ast.Return)
			Expect(ok).To(BeTrue())
		})
		It("can parse a return statement with a value", func() {
			r, ok := parse("return 42").(ast.Return)
			Expect(ok).To(BeTrue())
			expectNumber(r.Value(), 42)
		})
		It("can parse a when expression", func() {
			w, ok := parse("when (1) { 2 -> { 3 },, 4 -> { 5 }, else -> { 6 }, 7 -> { 8 }, else -> { 9 } }").(ast.When)
			Expect(ok).To(BeTrue())
			expectNumber(w.Target(), 1)
			Expect(len(w.Clauses())).To(Equal(5))
			wv, ok := w.Clauses()[0].(ast.WhenValueClause)
			Expect(ok).To(BeTrue())
			expectNumber(wv.Value(), 2)
			expectNumber(wv.Body(), 3)
			wv, ok = w.Clauses()[1].(ast.WhenValueClause)
			Expect(ok).To(BeTrue())
			expectNumber(wv.Value(), 4)
			expectNumber(wv.Body(), 5)
			we, ok := w.Clauses()[2].(ast.WhenElseClause)
			Expect(ok).To(BeTrue())
			expectNumber(we.Body(), 6)
			wv, ok = w.Clauses()[3].(ast.WhenValueClause)
			Expect(ok).To(BeTrue())
			expectNumber(wv.Value(), 7)
			expectNumber(wv.Body(), 8)
			we, ok = w.Clauses()[4].(ast.WhenElseClause)
			Expect(ok).To(BeTrue())
			expectNumber(we.Body(), 9)
		})
		It("can parse a when with boolean expressions", func() {
			parse(`...dyego
                when {
                  a > b -> { break }
                }
            `)
		})
	})
	Describe("types", func() {
		t := func(text string) ast.Element {
			v, ok := parse("var a: " + text).(ast.VarDefinition)
			Expect(ok).To(BeTrue())
			return v.Type()
		}
		Describe("expressions", func() {
			It("can parse a simple reference", func() {
				ty := t("A")
				n(ty, "A")
			})
			It("can parse a dotted reference", func() {
				ty := t("a.B")
				s, ok := ty.(ast.Selection)
				Expect(ok).To(BeTrue())
				n(s.Target(), "a")
				n(s.Member(), "B")
			})
			It("can parse a sequence type referenece", func() {
				ty := t("A[]")
				st, ok := ty.(ast.SequenceType)
				Expect(ok).To(BeTrue())
				n(st.Elements(), "A")
			})
			It("can parse an optional type reference", func() {
				ty := t("A?")
				ot, ok := ty.(ast.OptionalType)
				Expect(ok).To(BeTrue())
				n(ot.Element(), "A")
			})
		})
		Describe("type literal", func() {
			tl := func(text string) ast.TypeLiteral {
				r, ok := t(text).(ast.TypeLiteral)
				Expect(ok).To(BeTrue())
				return r
			}
			tm := func(e ast.Element) ast.TypeLiteralMember {
				r, ok := e.(ast.TypeLiteralMember)
				Expect(ok).To(BeTrue())
				return r
			}
			cm := func(e ast.Element) ast.CallableTypeMember {
				r, ok := e.(ast.CallableTypeMember)
				Expect(ok).To(BeTrue())
				return r
			}
			It("can parse an empty type", func() {
				ty := tl("<>")
				Expect(len(ty.Members())).To(Equal(0))
			})
			It("can parse a type with a member", func() {
				ty := tl("< a: Int >")
				m := tm(ty.Members()[0])
				n(m.Name(), "a")
				n(m.Type(), "Int")
			})
			It("can parse a type with two members", func() {
				ty := tl("< a: Int, b: Int >")
				m := tm(ty.Members()[0])
				n(m.Name(), "a")
				n(m.Type(), "Int")
				m = tm(ty.Members()[1])
				n(m.Name(), "b")
				n(m.Type(), "Int")
			})
			It("can parse a nested type literal", func() {
				ty := tl("<a: <a: Int>>")
				m := tm(ty.Members()[0])
				n(m.Name(), "a")
				nl := m.Type().(ast.TypeLiteral)
				m = tm(nl.Members()[0])
				n(m.Name(), "a")
				n(m.Type(), "Int")
			})
			It("can parse a callable member", func() {
				ty := tl("< { x: Int, y: Int -> Int } >")
				c := cm(ty.Members()[0])
				Expect(len(c.Parameters())).To(Equal(2))
				n(c.Result(), "Int")
			})
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
			It("can parse a precedence relation", func() {
				o := op("infix operator a before infix b left")
				p := o.Precedence()
				Expect(p.Relation()).To(Equal(ast.Before))
				Expect(p.Placement()).To(Equal(ast.Infix))
				o = op("infix operator a after prefix b right")
				p = o.Precedence()
				Expect(p.Relation()).To(Equal(ast.After))
				Expect(p.Placement()).To(Equal(ast.Prefix))
				o = op("infix operator a after postfix b right")
				p = o.Precedence()
				Expect(p.Relation()).To(Equal(ast.After))
				Expect(p.Placement()).To(Equal(ast.Postfix))
			})
		})
	})
	e := func(source string) ast.Element {
		actualSource := "...dyego," + source
		element := parse(actualSource)
		sequence := element.(ast.Sequence)
		return sequence.Right()
	}
	expectName := func(n ast.Element, name string) {
		nameElement, ok := n.(ast.Name)
		Expect(ok).To(BeTrue())
		Expect(nameElement.Text()).To(Equal(name))
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
	Describe("prefix expressions", func() {
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
		It("can parse an local identifier operator", func() {
			v := e("this dot this")
			l, r := expectBinaryOp(v, "dot")
			n(l, "this")
			n(r, "this")
		})
		It("can handle an assignment", func() {
			e(`
                when {
                    ret.hit -> {
                        hitSphere = obj
                        isHit = true
                        tval = hit.tval
                    }
                }
            `)
		})
	})
	Describe("separators", func() {
		sequence := func(e ast.Element) []ast.Element {
			var result []ast.Element
			s, ok := e.(ast.Sequence)
			for ok {
				l := s.Left()
				result = append(result, l)
				e = s.Right()
				s, ok = e.(ast.Sequence)
			}
			result = append(result, e)
			return result
		}
		s := func(text string) []ast.Element {
			e := parse(text)
			return sequence(e)
		}
		ns := func(e []ast.Element, names ...string) {
			Expect(len(e)).To(Equal(len(names)))
			for i, name := range names {
				n(e[i], name)
			}
		}
		It("new lines can imply sperators", func() {
			ns(s("a \n b"), "a", "b")
		})
		It("operator before a nl prevents implied separater", func() {
			seq := s("...dyego \n  a + \n b")
			Expect(len(seq)).To(Equal(2))
		})
		It("operator after a nl prevents implied separater", func() {
			seq := s("...dyego, a \n + b")
			Expect(len(seq)).To(Equal(2))
		})
	})
	Describe("locals", func() {
		dec := func(text string) ast.VarDefinition {
			v, ok := parse(text).(ast.VarDefinition)
			Expect(ok).To(BeTrue())
			return v
		}
		It("can parse a simple val", func() {
			v := dec("val a = 1")
			n(v.Name(), "a")
			i(v.Value(), 1)
			Expect(v.Type()).To(BeNil())
			Expect(v.Mutable()).To(BeFalse())
		})
		It("can parse a simple typed val", func() {
			v := dec("val a: Int = 1")
			n(v.Name(), "a")
			i(v.Value(), 1)
			n(v.Type(), "Int")
			Expect(v.Mutable()).To(BeFalse())
		})
		It("can parse a simple var", func() {
			v := dec("var a = 1")
			n(v.Name(), "a")
			i(v.Value(), 1)
			Expect(v.Type()).To(BeNil())
			Expect(v.Mutable()).To(BeTrue())
		})
		It("can parse a simple typed var", func() {
			v := dec("var a: Int = 1")
			n(v.Name(), "a")
			i(v.Value(), 1)
			n(v.Type(), "Int")
			Expect(v.Mutable()).To(BeTrue())
		})
		It("can parse an late initialized var", func() {
			v := dec("var a")
			n(v.Name(), "a")
			Expect(v.Value()).To(BeNil())
			Expect(v.Type()).To(BeNil())
			Expect(v.Mutable()).To(BeTrue())
		})
	})
	Describe("errors", func() {
		It("reports invalid vocabulary references", func() {
			expectErrors("...missing", "Expected a vocabulary reference")
		})
		It("reports an invalid expression", func() {
			expectErrors("(val a)", "Expected one of")
		})
		It("reports an invalid vocabulary member", func() {
			expectErrors("let a = <| invalid operator |>", "Expected one of infix")
		})
		It("reports an invalid let", func() {
			expectErrors("let a = <!", "Expected one of <|")
			expectErrors("let a = else", "Expected one of <|")
		})
		It("reports an invalid parameter", func() {
			expectErrors("...dyego, a(a: 1 + 2 + var)", "received var")
		})
		It("reports an invalid sequence", func() {
			expectErrors(".", "received .")
		})
		It("reports an undefined operator", func() {
			expectErrors("&&& b", "Symbol '&&&'")
		})
		It("reports an invalid call", func() {
			expectErrors("a(:1)", "Expected <identifier>")
		})
	})
	Describe("examples", func() {
		It("can parse the simple example", func() {
			parseFile("../examples/Simple.dg")
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
  infix operator identifiers left,
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
	return append(result, len(source)+1)
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
	if line >= len(lineMap) {
		return lineMap[len(lineMap)-1]
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
			println(fmt.Sprintf("%d:%d: %s", line, col, err.Message()))
			if line < len(lineMap) {
				println(source[lineMap[line-1] : lineMap[line]-1])
			}
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

func expectErrors(text string, errors ...string) {
	p := NewParser(scan(text), defaultScope)
	p.Parse()
loop:
	for _, message := range errors {
		for _, err := range p.Errors() {
			if strings.Contains(err.Message(), message) {
				continue loop
			}
		}
		printErrors(p.Errors(), text)
		Fail(fmt.Sprintf("Expected '%s' to be included as an error", message))
	}
}

func expectNumber(element ast.Element, value int) {
	n, ok := element.(ast.LiteralInt)
	Expect(ok).To(Equal(true))
	Expect(n.Value()).To(Equal(value))
}

func expectName(element ast.Element, value string) {
	l, ok := element.(ast.Name)
	Expect(ok).To(Equal(true))
	Expect(l.Text()).To(Equal(value))
}

func readFile(name string) []byte {
	src, err := ioutil.ReadFile(name)
	Expect(err).To(BeNil())
	return append(src, 0)
}

func parseFile(name string) ast.Element {
	return parse(string(readFile(name)))
}

func TestErrors(t *testing.T) {
	dyegoVocabulary := parseVocabulary(dyego0VocabularySource)
	scope := newVocabularyScope()
	scope.members["dyego"] = dyegoVocabulary
	defaultScope = scope

	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}
