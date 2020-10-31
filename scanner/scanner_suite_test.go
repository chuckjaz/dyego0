package scanner_test

import (
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego0/location"
	"dyego0/scanner"
	"dyego0/tokens"
)

func scanBytes(src []byte, expected ...tokens.Token) int {
	scanner := scanner.NewScanner(src, 0, nil)
	var received tokens.Token
	for _, token := range expected {
		received = scanner.Next()
		Ω(received).Should(Equal(token))
	}
	token := scanner.Next()
	Ω(token).Should(Equal(tokens.EOF))
	return scanner.Line()
}

func scanString(text string, expected ...tokens.Token) int {
	src := append([]byte(text), 0)
	return scanBytes(src, expected...)
}

func scannerOf(text string) *scanner.Scanner {
	src := append([]byte(text), 0)
	return scanner.NewScanner(src, 0, nil)
}

var _ = Describe("scanner", func() {
	Describe("when constructing the instance", func() {
		It("should not panic", func() {
			scanner.NewScanner([]byte{0}, 0, nil)
		})
	})
	It("should panic when an empty buffer is passed", func() {
		Ω(func() {
			scanner.NewScanner([]byte{}, 0, nil)
		}).Should(Panic())
	})
	It("should panic when a non-null-terminated buffer is passed", func() {
		Ω(func() {
			scanner.NewScanner([]byte{'a', 'b', 'c'}, 0, nil)
		})
	})
	Describe("when parsing", func() {
		It("should parse 'ident' as an IDENT", func() {
			scanString("ident", tokens.Identifier)
		})
		It("should parse multiple idents as IDENT", func() {
			scanString("  ident   ident2 _ _12", tokens.Identifier, tokens.Identifier, tokens.Identifier, tokens.Identifier)
		})
		It("should report start/end correctly", func() {
			scanString(" ident // comment", tokens.Identifier)
		})
		It("should count lines correctly", func() {
			lines := scanString(" \n  \r\n  \n ident ", tokens.Identifier)
			Ω(lines).Should(Equal(4))
		})
		It("should be able to clone a scanner", func() {
			s := scanner.NewScanner(append([]byte(" a b c "), 0), 0, nil)
			s.Next()
			c := s.Clone()
			Expect(s.Start()).To(Equal(location.Pos(1)))
			Expect(s.End()).To(Equal(location.Pos(2)))
			Expect(s.Start()).To(Equal(c.Start()))
			Expect(s.End()).To(Equal(c.End()))
		})
		It("can scan an integer qualifier", func() {
			scanString("1i", tokens.LiteralInt)
		})
		It("can scan a float qualifier", func() {
			scanString("1f", tokens.LiteralFloat)
		})
		It("report an invalid float", func() {
			scanString(strings.Repeat("9", 1000)+"f", tokens.Invalid)
		})
		It("report an invalid double", func() {
			scanString(strings.Repeat("9", 1000)+"d", tokens.Invalid)
		})
		It("can scan a double qualifier", func() {
			scanString("1d 1.0d", tokens.LiteralDouble, tokens.LiteralDouble)
		})
		It("can scan special runes", func() {
			scanString("'\\0' '\\n' '\\r' '\\b' '\\t' '\\\\' '\\'' ''", tokens.LiteralRune, tokens.LiteralRune,
				tokens.LiteralRune, tokens.LiteralRune, tokens.LiteralRune, tokens.LiteralRune, tokens.LiteralRune,
				tokens.Invalid)
		})
		It("can special runes in a string", func() {
			scanString("\" \\\" \\r \\b \\\\ \"", tokens.LiteralString)
		})
		It("can scan a escaped identifier", func() {
			scanString(" `+` ", tokens.Identifier)
		})
		It("can report an invalid escaped identifier", func() {
			scanString(" `  \n", tokens.Invalid)
		})
		It("can scan a integer range", func() {
			scanString("1..4", tokens.LiteralInt, tokens.Symbol, tokens.LiteralInt)
		})
		It("can indicate where the nl was", func() {
			s := scannerOf("a b \n c d \n e")
			Expect(s.NewLineLocation().IsValid()).To(BeFalse())
			s.Next()
			Expect(s.NewLineLocation().IsValid()).To(BeFalse())
			s.Next()
			Expect(s.NewLineLocation().IsValid()).To(BeFalse())
			s.Next()
			Expect(s.NewLineLocation()).To(Equal(location.Pos(4)))
			s.Next()
			Expect(s.NewLineLocation().IsValid()).To(BeFalse())
			s.Next()
			Expect(s.NewLineLocation()).To(Equal(location.Pos(10)))
		})
	})
})

func TestScanner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scanner test suite")
}
