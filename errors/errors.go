package errors

import (
	"dyego0/location"
	"fmt"
)

// Error is an error that is associated with a Pos range
type Error interface {
	error
	location.Locatable
}

// New returns a new formatted error message for the given locatable
func New(loc location.Locatable, message string, args ...interface{}) Error {
	return NewAt(loc.Start(), loc.End(), message, args...)
}

// NewAt returns a new formatted error located at start and ending at end
func NewAt(start location.Pos, end location.Pos, message string, args ...interface{}) Error {
	return &errorImpl{Location: location.NewLocation(start, end), message: fmt.Sprintf(message, args...)}
}

type errorImpl struct {
	location.Location
	message string
}

func (e *errorImpl) Error() string {
	return e.message
}
