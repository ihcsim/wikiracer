package wikiracer

import "fmt"

const (
	ErrPrefixDestinationUnreachable = "Destination unreachable"
	ErrPrefixPageNotFound           = "Page not found"
	ErrPrefixInvalidEmptyInputs     = "The provided inputs must not be empty"
)

// DestinationUnreachable is the error used when the destination isn't reachable on a given path.
type DestinationUnreachable struct {
	destination string
}

// Error returns the string representation of the DestinationUnreachable error.
func (e DestinationUnreachable) Error() string {
	return fmt.Sprintf("%s: %s", ErrPrefixDestinationUnreachable, e.destination)
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
