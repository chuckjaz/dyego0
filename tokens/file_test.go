package tokens_test

import (
	"fmt"
	"strconv"

	"dyego0/tokens"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("file", func() {
	It("can report an invalid position", func() {
		Expect(tokens.Pos(-1).IsValid()).To(BeFalse())
	})
	It("can report a valid position", func() {
		Expect(tokens.Pos(100).IsValid()).To(BeTrue())
	})
	It("can create a file set", func() {
		fs := tokens.NewFileSet()
		Expect(fs).To(Not(BeNil()))
	})
	It("can build a file builder", func() {
		fs := tokens.NewFileSet()
		fb := fs.BuildFile("somefile", 100)
		Expect(fb).To(Not(BeNil()))
	})
	It("can declare lines of a file", func() {
		fs := tokens.NewFileSet()
		fb := fs.BuildFile("somefile", 100)
		fb.AddLine(0)
		fb.AddLine(10)
		fb.AddLine(20)
		f := fb.Build()
		p := f.Pos(15)
		Expect(p.IsValid()).To(BeTrue())
		Expect(f.Column(p)).To(Equal(6))
		Expect(f.Line(p)).To(Equal(2))
	})
	It("can declare lines out of order", func() {
		fs := tokens.NewFileSet()
		fb := fs.BuildFile("somefile", 1000)
		fb.AddLine(0)
		fb.AddLine(20)
		fb.AddLine(40)
		fb.AddLine(10)
		fb.AddLine(30)
		f := fb.Build()
		p := f.Pos(15)
		Expect(p.IsValid()).To(BeTrue())
		Expect(f.Column(p)).To(Equal(6))
		Expect(f.Line(p)).To(Equal(2))
	})
	It("can declare lines out of order", func() {
		fs := tokens.NewFileSet()
		fb := fs.BuildFile("somefile", 1000)
		Expect(fb.Pos(50).IsValid()).To(BeTrue())
		Expect(fb.Pos(-10).IsValid()).To(BeFalse())
	})
	It("can convert a Pos to a Position", func() {
		fs := tokens.NewFileSet()
		var files []tokens.File
		var positions []tokens.Pos
		for i := 0; i < 10; i++ {
			fb := fs.BuildFile("somefile"+strconv.Itoa(i), 1000)
			for l := 0; l < 10; l++ {
				fb.AddLine(l * 80)
			}
			f := fb.Build()
			files = append(files, f)
			positions = append(positions, f.Pos(150))
		}
		for i := 0; i < 10; i++ {
			p := positions[i]
			Expect(p.IsValid()).To(BeTrue())
			position := fs.Position(positions[i])
			Expect(position.IsValid()).To(BeTrue())
			Expect(position.FileName()).To(Equal("somefile" + strconv.Itoa(i)))
			Expect(position.Column()).To(Equal(71))
			Expect(position.Line()).To(Equal(2))
			Expect(position.String()).To(Equal(fmt.Sprintf("somefile%d:%d:%d", i, position.Line(), position.Column())))
		}
	})
	It("can build files out of order", func() {
		fs := tokens.NewFileSet()
		var builders []tokens.FileBuilder
		files := make([]tokens.File, 10)
		positions := make([]tokens.Pos, 10)
		for i := 0; i < 10; i++ {
			fb := fs.BuildFile("somefile"+strconv.Itoa(i), 1000)
			for l := 0; l < 10; l++ {
				fb.AddLine(l * 80)
			}
			builders = append(builders, fb)
		}
		for i := 0; i < 10; i += 2 {
			f := builders[i].Build()
			files[i] = f
			positions[i] = f.Pos(150)
		}
		for i := 1; i < 10; i += 2 {
			f := builders[i].Build()
			files[i] = f
			positions[i] = f.Pos(150)
		}
		for i := 0; i < 10; i++ {
			p := positions[i]
			Expect(p.IsValid()).To(BeTrue())
			position := fs.Position(positions[i])
			Expect(position.IsValid()).To(BeTrue())
			Expect(position.FileName()).To(Equal("somefile" + strconv.Itoa(i)))
			Expect(position.Column()).To(Equal(71))
			Expect(position.Line()).To(Equal(2))
			Expect(position.String()).To(Equal(fmt.Sprintf("somefile%d:%d:%d", i, position.Line(), position.Column())))
		}
	})
	It("can report an out of bound file", func() {
		fs := tokens.NewFileSet()
		Expect(fs.Position(tokens.Pos(10)).IsValid()).To(BeFalse())
		fb := fs.BuildFile("somefile", 1000)
		fb.Build()
		position := fs.Position(tokens.Pos(-1))
		Expect(position.IsValid()).To(BeFalse())
		Expect(position.String()).To(Equal("invalid"))
	})
	It("can report an out of bound lines", func() {
		fs := tokens.NewFileSet()
		fb := fs.BuildFile("somefile", 1000)
		fb.AddLine(10)
		fb.AddLine(20)
		f := fb.Build()
		Expect(f.Line(tokens.Pos(-100))).To(Equal(-1))
		Expect(f.Line(tokens.Pos(2000))).To(Equal(-1))
		Expect(f.Line(tokens.Pos(0))).To(Equal(1))
	})
	It("can report an out of bound columns", func() {
		fs := tokens.NewFileSet()
		fb := fs.BuildFile("somefile", 1000)
		fb.AddLine(0)
		fb.AddLine(10)
		fb.AddLine(20)
		f := fb.Build()
		Expect(f.Column(tokens.Pos(-100))).To(Equal(-1))
		Expect(f.Column(tokens.Pos(2000))).To(Equal(-1))
		Expect(f.Column(tokens.Pos(0))).To(Equal(1))
	})
	It("can convert to an offset", func() {
		fs := tokens.NewFileSet()
		f0 := fs.BuildFile("somefile0", 1000).Build()
		f1 := fs.BuildFile("somefile1", 1000).Build()
		Expect(f0.Offset(tokens.Pos(500))).To(Equal(500))
		Expect(f1.Offset(tokens.Pos(1500))).To(Equal(500))
	})
	It("can get a pos from a file", func() {
		fs := tokens.NewFileSet()
		f0 := fs.BuildFile("somefile0", 1000).Build()
		f1 := fs.BuildFile("somefile0", 1000).Build()
		Expect(f0.Pos(500)).To(Equal(tokens.Pos(500)))
		Expect(f1.Pos(500)).To(Equal(tokens.Pos(1500)))
		Expect(f0.Pos(1500)).To(Equal(tokens.Pos(-1)))
		Expect(f1.Pos(-1)).To(Equal(tokens.Pos(-1)))
	})
	It("can get a file size", func() {
		fs := tokens.NewFileSet()
		f0 := fs.BuildFile("somefile0", 1000).Build()
		f1 := fs.BuildFile("somefile0", 1000).Build()
		Expect(f0.Size()).To(Equal(1000))
		Expect(f1.Size()).To(Equal(1000))
	})
})
