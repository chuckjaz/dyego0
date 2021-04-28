package sources_test

import (
	"dyego0/sources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("sourcefilereaderss", func() {
	emptyReader := func(fileName string) (io.Reader, error) {
		return nil, nil
	}
	It("can create a file mdule scope", func() {
		scope, err := sources.NewSourceFileReaderScope(nil, emptyReader)
		Expect(scope).To(Not(BeNil()))
		Expect(err).To(BeNil())
	})
	It("can create and use a scope with files", func() {
		root, err := sources.NewSourceFileReaderScope([]string{
			"a/b/c.go",
			"a/b/d.go",
			"a/e/g.go",
			"h/i/j/k.go",
			"h/i/l/m.go",
		}, emptyReader)
		Expect(err).To(BeNil())
		a, err := root.FindScope("a")
		Expect(err).To(BeNil())
		Expect(a).To(Not(BeNil()))
		b, err := a.FindScope("b")
		Expect(b).To(Not(BeNil()))
		Expect(err).To(BeNil())
		c, err := b.Find("c")
		Expect(c).To(Not(BeNil()))
		Expect(err).To(BeNil())
		Expect(c.Name()).To(Equal("c"))
		Expect(c.FileName()).To(Equal("a/b/c.go"))
		r, err := c.NewReader()
		Expect(r).To(BeNil())
		Expect(err).To(BeNil())
		_, err = root.Find("a")
		Expect(err).To(Not(BeNil()))
		_, err = root.FindScope("missing")
		Expect(err).To(Not(BeNil()))
	})
	It("reports duplicate files", func() {
		_, err := sources.NewSourceFileReaderScope([]string{
			"a/b/c.go",
			"a/b/c.go",
		}, emptyReader)
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(Equal("Duplicate file 'a/b/c.go'"))
	})
	It("reports invalid file", func() {
		_, err := sources.NewSourceFileReaderScope([]string{
			"",
		}, emptyReader)
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(Equal("Invalid file name ''"))
	})
})
