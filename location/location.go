package location

import (
	"fmt"
)

// Pos is a unique source location information
type Pos int

// IsValid returns if the position is valid. Invalid Pos values are returned, for
// example, if the request offset negative or past the end of the file.
func (p Pos) IsValid() bool {
	return p >= 0
}

// Locatable allows locating an item in a source file
type Locatable interface {
	Start() Pos
	End() Pos
	Length() int
}

// Location is the location of an item in a source file
type Location struct {
	start, end Pos
}

// NewLocation creates a new location
func NewLocation(start, end Pos) Location {
	return Location{start, end}
}

// Start is the start of the item in a source file
func (l Location) Start() Pos {
	return l.start
}

// End is the end of the item in a source file
func (l Location) End() Pos {
	return l.end
}

// Length is the length of the item (End() - Start().
func (l Location) Length() int {
	return int(l.end) - int(l.start)
}

func (l Location) String() string {
	return fmt.Sprintf("Location(%d-%d)", l.start, l.end)
}
