package location_test

import (
    "testing"
	"dyego0/location"
	"dyego0/tokens"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("location", func() {
	s := tokens.Pos(1)
	e := tokens.Pos(100)
	l := location.NewLocation(s, e)

	It("should access start", func() {
		Expect(l.Start()).To(Equal(s))
	})

	It("should access end", func() {
		Expect(l.End()).To(Equal(e))
	})

	It("should access length", func() {
		Expect(l.Length()).To(Equal(99))
	})
})

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Location Suite")
}

