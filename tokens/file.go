package tokens

import (
	"fmt"
	"sort"

	"dyego0/location"
)

// File access line column inforamtion for a source file
type File interface {
	// Column is the 1-based column of the given Pos the file
	Column(p location.Pos) int

	// FileName is the file name given when declared in the FileSet
	FileName() string

	// Line is the 1-based line of the given Pos in the file
	Line(pos location.Pos) int

	// LineStart returns the 0 based offset of the start of the 1 based line. If the line is
	// past the last line the offset is the size of the file.
	LineStart(line int) int

	// Offset is the 0-based offset into the file using UTF-8 encoding
	Offset(p location.Pos) int

	// Calculate a Pos for the given 0-based offset using UTF-8 encoding
	Pos(offset int) location.Pos

	// Return a Position for the given Pos in the file. Position is an easier to use
	// encoding of file position but is signficant larger than a Pos which is specialized
	// int
	Position(p location.Pos) Position

	// Give the size of the source file in UTF-8 encoding
	Size() int
}

// FileBuilder is returned by a FileSet to allow building the definition of a File.
type FileBuilder interface {
	// AddLine declares the offset of a line. Lines can be declared in any order but it
	// is more efficient to declare them in order.
	AddLine(offset int)

	// Pos calculates a Pos for the given 0-based offset using UTF-8 encoding.
	Pos(offset int) location.Pos

	// Build finalizes the definition of the File and returns the immutable File defined.
	Build() File
}

// FileSet is a set of File defintions for an arbitrary number of files. All Pos values
// for a FileSet are unique to the File in the set.
type FileSet interface {
	// BuildFile declares a file in the file set. It returns a FileBuilder which allows
	// building the File defintiion which is immutable after it is built.
	BuildFile(filename string, size int) FileBuilder

	// Position returns a Position which gives access allows calculating the Line, Column
	// and FileName of Pos for any file in the FileSet. Pos are only 4 bytes and allow
	// efficient encoding of file/line/column information. The Position is significantly
	// larger encoding of the same information.
	Position(p location.Pos) Position

	// File returns the File associated with pos. If pos is invalid returns nil.
	File(pos location.Pos) File
}

// Position is an easier to use, but signficantly larger, encoding of a Pos that allows
// easy access to the FileName/Line/Column information.
type Position interface {
	// FileName is the name of the file declared with FileSet.BuildFile
	FileName() string

	// Column is a 1-based column of the source position which is the number of UTF-8
	// code units between the source position and the start of the line.
	Column() int

	// Line is a 1-based line of the soruce position.
	Line() int

	// IsValid returns if the underlying position was in the file requested or was itself
	// valid.
	IsValid() bool

	// String return the position in <filename>:<line>:<column> format which can be used
	// in displaying a diagnostic message for a source location.
	String() string
}

// NewFileSet returns a FileSet instance that allows declaring a set of files allowing
// efficient encoding of source location infomration in Pos values that unique identify
// the location in a set of source files in 4 bytes.
func NewFileSet() FileSet {
	return &fileSet{}
}

type lines []int

func (ls lines) Search(x int) int {
	return sort.SearchInts(ls, x)
}

func (ls lines) Insert(index, value int) lines {
	result := append(ls, 0)
	copy(result[index+1:], result[index:])
	result[index] = value
	return result
}

type fileBuilder struct {
	filename string
	base     int
	size     int
	lines    lines
	fileSet  *fileSet
}

func (fb *fileBuilder) AddLine(offset int) {
	l := len(fb.lines)
	if l == 0 || fb.lines[l-1] < offset {
		fb.lines = append(fb.lines, offset)
	} else {
		index := fb.lines.Search(offset)
		if fb.lines[index] != offset {
			fb.lines = fb.lines.Insert(index, offset)
		}
	}
}

func (fb *fileBuilder) Pos(offset int) location.Pos {
	if offset < 0 || offset > fb.size {
		return location.Pos(-1)
	}
	return location.Pos(fb.base + offset)
}

func (fb *fileBuilder) Build() File {
	result := &file{filename: fb.filename, base: fb.base, size: fb.size, lines: fb.lines}
	fb.fileSet.add(result)
	return result
}

type fileSet struct {
	base  int
	files files
}

func (fs *fileSet) add(file *file) {
	fs.files = fs.files.add(file)
}

func (fs *fileSet) BuildFile(filename string, size int) FileBuilder {
	b := fs.base
	fs.base = b + size
	return &fileBuilder{filename: filename, size: size, base: b, fileSet: fs, lines: lines{0}}
}

func (fs *fileSet) File(pos location.Pos) File {
	l := len(fs.files)
	if l == 0 {
		return nil
	}
	if !pos.IsValid() {
		return nil
	}
	index := fs.files.Search(int(pos))
	if index == l {
		index--
	}
	f := fs.files[index]
	if f.base > int(pos) && index > 0 {
		f = fs.files[index-1]
	}
	return f
}

func (fs *fileSet) Position(pos location.Pos) Position {
	f := fs.File(pos)
	if f == nil {
		return &position{pos: location.Pos(-1)}
	}
	return f.Position(pos)
}

type file struct {
	filename string
	base     int
	size     int
	lines    lines
}

func (f *file) lineOf(p location.Pos) int {
	o := int(p) - f.base
	if o < 0 || o > f.size {
		return -2
	}
	line := f.lines.Search(o)
	if line == len(f.lines) && o < f.size {
		return len(f.lines) - 1
	}
	if line >= len(f.lines) {
		return len(f.lines)
	}
	ls := f.lines[line]
	if ls > o {
		if line > 0 {
			line--
		} else {
			line = 0
		}
	}
	return line
}

func (f *file) linePos(p location.Pos) location.Pos {
	line := f.lineOf(p)
	if line < 0 {
		return location.Pos(-1)
	}
	if line >= len(f.lines) {
		return location.Pos(f.size)
	}
	return location.Pos(f.lines[line] + f.base)
}

func (f *file) Column(p location.Pos) int {
	lp := f.linePos(p)
	if int(lp) < 0 {
		return -1
	}
	return int(p-lp) + 1
}

func (f *file) FileName() string {
	return f.filename
}

func (f *file) Line(p location.Pos) int {
	return f.lineOf(p) + 1
}

func (f *file) LineStart(line int) int {
	offset := line - 1
	if offset < 0 {
		return 0
	}
	if offset >= len(f.lines) {
		return f.size
	}
	return f.lines[offset]
}

func (f *file) Offset(p location.Pos) int {
	return int(p) - f.base
}

func (f *file) Pos(offset int) location.Pos {
	if offset < 0 || offset > f.size {
		return location.Pos(-1)
	}
	return location.Pos(f.base + offset)
}

func (f *file) Position(p location.Pos) Position {
	o := int(p) - f.base
	if o < 0 || o > f.size {
		return &position{file: f, pos: location.Pos(-1)}
	}
	return &position{file: f, pos: p}
}

func (f *file) Size() int {
	return f.size
}

type files []*file

func (fs files) Search(x int) int {
	return sort.Search(len(fs), func(index int) bool {
		return fs[index].base >= x
	})
}

func (fs files) add(file *file) files {
	l := len(fs)
	if l == 0 || fs[l-1].base < file.base {
		return append(fs, file)
	}
	index := fs.Search(file.base)
	if fs[index] != file {
		fs = fs.Insert(index, file)
	}
	return fs
}

func (fs files) Insert(index int, file *file) files {
	result := append(fs, nil)
	copy(result[index+1:], result[index:])
	result[index] = file
	return result
}

type position struct {
	file *file
	pos  location.Pos
}

func (p *position) FileName() string {
	return p.file.FileName()
}

func (p *position) Column() int {
	return p.file.Column(p.pos)
}

func (p *position) Line() int {
	return p.file.Line(p.pos)
}

func (p *position) IsValid() bool {
	return p.pos.IsValid()
}

func (p *position) String() string {
	if p.pos.IsValid() {
		line := p.Line()
		column := p.Column()
		filename := p.FileName()
		if filename != "" && line >= 0 && column >= 0 {
			return fmt.Sprintf("%s:%d:%d", filename, line, column)
		}
	}
	return "invalid"
}
