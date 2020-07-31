package parser

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "dyego0/ast"
    "dyego0/scanner"
)

var _ = Describe("table", func() {
    Describe("prcedenceLebel" , func() {
        It("can create a precedence level", func() {
            level := newPrecedenceLevel()
            Expect(level).To(Not(BeNil()))
        })
        It("can create a lower precedence", func() {
            level := newPrecedenceLevel()
            lower := level.MakeLower()
            Expect(lower).To(Not(BeNil()))
            Expect(level).To(Equal(lower.Higher()))
            Expect(lower).To(Equal(level.Lower()))
            Expect(lower.IsHigherThan(level)).To(Equal(false))
            Expect(level.IsHigherThan(lower)).To(Equal(true))
        })
        It("can create a higher precedence", func() {
            level := newPrecedenceLevel()
            higher := level.MakeHigher()
            Expect(higher).To(Not(BeNil()))
            Expect(level).To(Equal(higher.Lower()))
            Expect(higher).To(Equal(level.Higher()))
            Expect(higher.IsHigherThan(level)).To(Equal(true))
            Expect(level.IsHigherThan(higher)).To(Equal(false))
        })
        It("transitive compare of lower", func() {
            a := newPrecedenceLevel()
            b := a.MakeLower()
            c := b.MakeLower()
            Expect(a.IsHigherThan(c)).To(Equal(true))
        })
        It("transitive compare of higher", func() {
            a := newPrecedenceLevel()
            b := a.MakeHigher()
            c := b.MakeHigher()
            Expect(c.IsHigherThan(a)).To(Equal(true))
        })
        It("out of order higher", func() {
            a := newPrecedenceLevel()
            c := a.MakeLower()
            b := c.MakeHigher()
            Expect(a.IsHigherThan(c)).To(Equal(true))
            Expect(a.IsHigherThan(b)).To(Equal(true))
        })
        It("out of order lower", func() {
            a := newPrecedenceLevel()
            c := a.MakeHigher()
            b := c.MakeLower()
            Expect(c.IsHigherThan(a)).To(Equal(true))
            Expect(c.IsHigherThan(b)).To(Equal(true))
        })
        It("unrelated", func() {
            a := newPrecedenceLevel()
            b := newPrecedenceLevel()
            Expect(b.IsHigherThan(a)).To(Equal(false))
            Expect(a.IsHigherThan(b)).To(Equal(false))
        })
        It("not higher than self", func() {
            a := newPrecedenceLevel()
            Expect(a.IsHigherThan(a)).To(Equal(false))
        })
    })
    Describe("operator", func() {
        It("can create a operator", func() {
            op := newOperator("+", nil, nil)
            Expect(op.Name()).To(Equal("+"))
            Expect(op.Levels()).To(BeNil())
            Expect(op.Associativities()).To(BeNil())
        })
    })
    Describe("vocabulary", func() {
        It("can create a vocabulary", func() {
            vocab := newVocabulary()
            Expect(vocab).To(Not(BeNil()))
            Expect(vocab.Scope()).To(Not(BeNil()))
        })
        It("can set and get a value", func() {
            vocab := newVocabulary()
            op := newOperator("+", nil, nil)
            vocab.members["+"] = op
            value, ok := vocab.Get("+")
            Expect(ok).To(Equal(true))
            Expect(value).To(Equal(op))
        })
    })
    Describe("vocabularyScope", func() {
        It("can crate a vocabulary scope", func() {
            scope := newVocabularyScope()
            Expect(scope).To(Not(BeNil()))
        })
        It("can get a scope value", func() {
            scope := newVocabularyScope()
            vocab := newVocabulary()
            scope.members["a"] = vocab
            result, ok := scope.Get("a")
            Expect(ok).To(Equal(true))
            Expect(result).To(Equal(vocab))
        })
    })
    Describe("buildVocabulary", func() {
        scan := func(text string) *scanner.Scanner {
            return scanner.NewScanner(append([]byte(text), 0), 0)
        }
        b := ast.NewBuilder(scan(""))
        b.PushContext()

        var scope vocabularyScope = newVocabularyScope()

        build := func(lit ast.VocabularyLiteral) vocabulary {
            result, errors := buildVocabulary(scope, lit)
            if len(errors) > 0 {
                Expect(len(errors)).To(Equal(0))
            }
            return result
        }

        buildError := func(lit ast.VocabularyLiteral, expects ...string) {
            _, errors := buildVocabulary(scope, lit)
            for i, err := range errors {
                Expect(err.message).To(Equal(expects[i]))
            }
            Expect(len(errors)).To(Equal(len(expects)))
        }

        vl := func(members ...ast.Element) ast.VocabularyLiteral {
            return b.VocabularyLiteral(members)
        }

        ns := func(names []string) []ast.Name {
            var result []ast.Name
            for _, name := range names {
                result = append(result, b.Name(name))
            }
            return result
        }

        op := func(
            placement ast.OperatorPlacement,
            associativity ast.OperatorAssociativity,
            precedence ast.VocabularyOperatorPrecedence,
            names ...string,
        ) ast.VocabularyOperatorDeclaration {
            return b.VocabularyOperatorDeclaration(ns(names), placement, precedence, associativity)
        }

        clike := build(
            vl(
                op(ast.Prefix, ast.Left, nil, "+", "-"),
                op(ast.Infix, ast.Left, nil, "*", "/", "%"),
                op(ast.Infix, ast.Left, nil, "+", "-"),
                op(ast.Infix, ast.Left, nil, "&&"),
                op(ast.Infix, ast.Left, nil, "||"),
            ),
        )

        s := func(members ...mbr) vocabularyScope {
            result := newVocabularyScope()
            for _, member := range members {
                result.members[member.name] = member.value
            }
            return result
        }

        scope = s(m("c", clike), m("s1", s(m("s2", s(m("c", clike))))), m("invalid", 1))

        getOp := func(v vocabulary, name string) operator {
            result, ok := v.Get(name)
            Expect(ok).To(Equal(true))
            Expect(result).To(Not(BeNil()))
            return result.(operator)
        }

        It("can build an empty vocabulary", func() {
            v := build(vl())
            Expect(v).To(Not(BeNil()))
        })
        It("can define an operator", func() {
            v :=  build(vl(op(ast.Infix, ast.Left, nil, "+")))
            o := getOp(v, "+")
            Expect(o.Levels()[ast.Infix]).To(Not(BeNil()))
            Expect(o.Levels()[ast.Prefix]).To(BeNil())
            Expect(o.Levels()[ast.Postfix]).To(BeNil())
            Expect(o.Associativities()[ast.Infix]).To(Equal(ast.Left))
            Expect(o.Associativities()[ast.Prefix]).To(Equal(ast.UnspecifiedAssociativity))
            Expect(o.Associativities()[ast.Postfix]).To(Equal(ast.UnspecifiedAssociativity))
        })
        It("can define operators at the same precedence", func() {
            v := build(vl(op(ast.Infix, ast.Left, nil, "+", "-")))
            plus := getOp(v, "+")
            minus := getOp(v, "-")
            Expect(plus.Levels()[ast.Infix]).To(Equal(minus.Levels()[ast.Left]))
        })
        It("can define operators at different levels", func() {
            v := build(vl(op(ast.Infix, ast.Left, nil, "*"), op(ast.Infix, ast.Left, nil, "+")))
            mult := getOp(v, "*")
            plus := getOp(v, "+")
            Expect(mult.Levels()[ast.Infix].IsHigherThan(plus.Levels()[ast.Infix])).To(Equal(true))
        })

        embed := func(names ...string) ast.VocabularyEmbedding {
           return b.VocabularyEmbedding(ns(names))
        }

        ref := func(name string, placement ast.OperatorPlacement, relation ast.OperatorPrecedenceRelation) ast.VocabularyOperatorPrecedence {
            return b.VocabularyOperatorPrecedence(b.Name(name), placement, relation)
        }

        It("can embedd a vocabulary and declare operator after", func() {
            v := build(
                vl(
                    embed("c"),
                    op(ast.Infix, ast.Left, ref("+", ast.Infix, ast.After), "<<", ">>"),
                ),
            )
            plus := getOp(v, "+")
            shift := getOp(v, ">>")
            Expect(plus.Levels()[ast.Infix].IsHigherThan(shift.Levels()[ast.Infix])).To(Equal(true))
        })
        It("can embedd a vocabulary and declare operator after", func() {
            v := build(
                vl(
                    embed("c"),
                    op(ast.Infix, ast.Left, ref("+", ast.Infix, ast.Before), "<<", ">>"),
                ),
            )
            plus := getOp(v, "+")
            shift := getOp(v, ">>")
            Expect(shift.Levels()[ast.Infix].IsHigherThan(plus.Levels()[ast.Infix])).To(Equal(true))
        })
        It("can embedd a vocabulary and reference no placement", func() {
            v := build(
                vl(
                    embed("c"),
                    op(ast.Infix, ast.Left, ref("*", ast.UnspecifiedPlacement, ast.After), "<<", ">>"),
                ),
            )
            mult := getOp(v, "*")
            shift := getOp(v, ">>")
            Expect(mult.Levels()[ast.Infix].IsHigherThan(shift.Levels()[ast.Infix])).To(Equal(true))
        })
        It("undeclared embedding", func() {
            buildError(vl(embed("cpp")), "Undefined vocabulary 'cpp'")
        })
        It("expected as scope", func() {
            buildError(vl(embed("c", "nested")), "Expected 'c' to be a vocabulary scope")
        })
        It("expect a vocabulary", func() {
            buildError(vl(embed("s1")), "Expected 's1' to be a vocabulary")
        })
        It("operator already defined", func() {
            buildError(vl(embed("c"), op(ast.Infix, ast.Left, nil, "+")), "An infix operator '+' already defined")
        })
        It("undeclared identifier", func() {
            buildError(
                vl(embed("c"), op(ast.Infix, ast.Left, ref("<>", ast.Infix, ast.After), "!=")),
                "Undeclared identifier '<>'",
            )
        })
        It("ambigious reference", func() {
            buildError(
                vl(embed("c"), op(ast.Infix, ast.Left, ref("+", ast.UnspecifiedPlacement, ast.After), "!=")),
                "Ambigious operator reference, both infix and prefix are defined",
            )
        })
        It("no placement defined", func() {
            buildError(
                vl(embed("c"), op(ast.Infix, ast.Left, ref("+", ast.Postfix, ast.After), "!=")),
                "No postfix placement defined for operator '+'",
            )
        })
    })
})

type mbr struct {
    name string
    value any
}

func m(name string, value any) mbr {
    return mbr{name: name, value: value}
}
