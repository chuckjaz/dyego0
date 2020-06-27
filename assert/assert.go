package assert

import "fmt"

// Assert panics when value is false
func Assert(value bool, message string, a ...interface{}) {
	if !value {
		panic(fmt.Sprintf(message, a...))
	}
}

// Fail will panic for an unconditionally failure
func Fail(message string, a ...interface{}) {
	panic(fmt.Sprintf(message, a...))
}

// Debug uses fmt.Printf to emit data. Useful for debugging
func Debug(message string, a ...interface{}) {
	fmt.Printf(message, a...)
	fmt.Println()
}
