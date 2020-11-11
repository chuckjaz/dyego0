package binder_test

import (
	"dyego0/binder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("context", func() {
	It("should be able to create a context", func() {
		context := binder.NewContext()
		Expect(context).To(Not(BeNil()))
	})
})
