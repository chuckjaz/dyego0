package assert_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dyego0/assert"
)

var _ = Describe("Assert", func() {
	Describe("When using Assert", func() {
		It("should not panic on success", func() {
			Expect(func() {
				assert.Assert(true, "this is true")
			}).NotTo(Panic())
		})
		It("should panic on failure", func() {
			Expect(func() {
				assert.Assert(false, "this is not true")
			}).To(Panic())
		})
	})
	Describe("When using Assert parameters", func() {
		It("should not panic on success", func() {
			Expect(func() {
				assert.Assert(true, "this is %t", true)
			}).NotTo(Panic())
		})
		It("should panic on failure", func() {
			Expect(func() {
				assert.Assert(false, "this is %t", false)
			}).To(Panic())
		})
	})
	Describe("when using Fail", func() {
		It("should panic", func() {
			Expect(func() {
				assert.Fail("Failure")
			}).To(Panic())
		})
	})
	Describe("when using Debug", func() {
		It("should not panic", func() {
			assert.Debug("some message %d", 1)
		})
	})
})

func TestAsserts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Asserts Suite")
}
