package sources_test

import (
	"dyego0/sources"
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sources", func() {
	Describe("SourceReader", func() {
		It("can create a source reader", func() {
			sourceReader := sources.NewSourceReader("a", "b", func() (io.Reader, error) { return nil, nil })
			Expect(sourceReader.Name()).To(Equal("a"))
			Expect(sourceReader.FileName()).To(Equal("b"))
			reader, err := sourceReader.NewReader()
			Expect(reader).To(BeNil())
			Expect(err).To(BeNil())
		})
	})
})

func TestModules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sources Suite")
}
