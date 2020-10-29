package ast

import (
	"dyego0/tokens"
	"fmt"
)

// Locatable allows locating an item in a source file
type Locatable interface {
	Start() tokens.Pos
	End() tokens.Pos
	Length() int
}

// Location is the location of an item in a source file
type Location struct {
	start, end tokens.Pos
}

// NewLocation creates a new location
func NewLocation(start, end tokens.Pos) Location {
	return Location{start, end}
}

// Start is the start of the item in a source file
func (l Location) Start() tokens.Pos {
	return l.start
}

// End is the end of the item in a source file
func (l Location) End() tokens.Pos {
	return l.end
}

// Length is the length of the item (End() - Start().
func (l Location) Length() int {
	return int(l.end) - int(l.start)
}

func (l Location) String() string {
	return fmt.Sprintf("Location(%d-%d)", l.start, l.end)
}
