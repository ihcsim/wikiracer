package wikiracer

import "fmt"

const (
	ErrPrefixNoLinksFound       = "The two pages aren't connected"
	ErrPrefixPageNotFound       = "Page not found"
	ErrPrefixInvalidEmptyInputs = "The provided inputs must not be empty"
)

// NoLinksFound is the error used when an origin isn't connected to a destination.
type NoLinksFound struct {
	origin      string
	destination string
}

// Error returns the string representation of the NoLinksFound error.
func (e NoLinksFound) Error() string {
	return fmt.Sprintf("%s: (%s, %s)", ErrPrefixNoLinksFound, e.origin, e.destination)
}

// PageNotFound is the error used when a non-existent page is requested.
type PageNotFound struct {
	Page
}

// Error returns the string representation of the PageNotFound error.
func (e PageNotFound) Error() string {
	return fmt.Sprintf("%s: %s", ErrPrefixPageNotFound, e.Title)
}

// InvalidEmptyInput is the error used when the provided inputs are empty.
type InvalidEmptyInput struct {
	origin      string
	destination string
}

// Error returns the string representation of the InvalidEmptyInput error.
func (e InvalidEmptyInput) Error() string {
	return fmt.Sprintf("%s: (%s, %s)", ErrPrefixInvalidEmptyInputs, e.origin, e.destination)
}
