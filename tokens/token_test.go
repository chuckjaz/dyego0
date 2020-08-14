package tokens_test

import (
	"dyego0/tokens"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("tokens", func() {
	Describe("Token.String()", func() {
		It("should convert identifier correctly", func() {
			Expect(tokens.Identifier.String()).To(Equal("<identifier>"))
		})
		It("should report an invalid token as invalid", func() {
			Expect(tokens.Token(1e6).String()).To(Equal("<invalid>"))
		})
		It("should expect all valid tokens to have valid strings", func() {
			for token := tokens.Token(0); token < tokens.InvalidToken; token++ {
				Expect(len(token.String()) > 0).To(BeTrue())
			}
		})
	})
	Describe("PseudoToken.String()", func() {
		It("should convert left correctly", func() {
			Expect(tokens.Left.String()).To(Equal("left"))
		})
		It("should report an unknown pseudo token as invalid", func() {
			Expect(tokens.PseudoToken(1e6).String()).To(Equal("<invalid>"))
		})
		It("should give valid strings for valid psuedo tokens", func() {
			for pseudoToken := tokens.PseudoToken(0); pseudoToken < tokens.InvalidPseudoToken; pseudoToken++ {
				Expect(len(pseudoToken.String()) > 0).To(BeTrue())
			}
		})
	})
})

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}
