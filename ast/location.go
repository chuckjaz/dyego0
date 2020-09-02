package ast

import (
	"fmt"
	"go/token"
)

// Locatable allows locating an item in a source file
type Locatable interface {
	Start() token.Pos
	End() token.Pos
	Length() int
}

// Location is the location of an item in a source file
type Location struct {
	start, end token.Pos
}

// NewLocation creates a new location
func NewLocation(start, end token.Pos) Location {
	return Location{start, end}
}

// Start is the start of the item in a source file
func (l Location) Start() token.Pos {
	return l.start
}

// End is the end of the item in a source file
func (l Location) End() token.Pos {
	return l.end
}

// Length is the length of the item (End() - Start().
func (l Location) Length() int {
	return int(l.end) - int(l.start)
}

func (l Location) String() string {
	return fmt.Sprintf("Location(%d-%d)", l.start, l.end)
}
