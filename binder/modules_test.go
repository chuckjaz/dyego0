package binder_test

import (
	"dyego0/binder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"testing"
)

var _ = Describe("modules", func() {
	Describe("ModuleSource", func() {
		It("can create a mdoule source", func() {
			moduleSource := binder.NewModuleSource("a", "b", func() (io.Reader, error) { return nil, nil })
			Expect(moduleSource.Name()).To(Equal("a"))
			Expect(moduleSource.FileName()).To(Equal("b"))
			reader, err := moduleSource.NewReader()
			Expect(reader).To(BeNil())
			Expect(err).To(BeNil())
		})
	})
})

func TestModules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}
