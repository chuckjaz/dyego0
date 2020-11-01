package errors_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego0/errors"
	"dyego0/location"
)

var _ = Describe("errors", func() {
	It("can create a new error using NewAt", func() {
		err := errors.NewAt(location.Pos(1), location.Pos(10), "Message")
		Expect(err.Error()).To(Equal("Message"))
	})
	It("can create a new error using New", func() {
		l := location.NewLocation(location.Pos(1), location.Pos(10))
		err := errors.New(l, "Message %d %d", 1, 2)
		Expect(err.Error()).To(Equal("Message 1 2"))
	})
})

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ast Suite")
}
