package location_test

import (
	"dyego0/location"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var _ = Describe("location", func() {
	s := location.Pos(1)
	e := location.Pos(100)
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

	It("can report an invalid position", func() {
		Expect(location.Pos(-1).IsValid()).To(BeFalse())
	})
	It("can report a valid position", func() {
		Expect(location.Pos(100).IsValid()).To(BeTrue())
	})
})

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Location Suite")
}
