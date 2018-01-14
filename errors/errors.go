package errors

import (
	"fmt"

	"github.com/ihcsim/wikiracer/internal/wiki"
)

const (
	ErrPrefixDestinationUnreachable = "Destination unreachable"
	ErrPrefixPageNotFound           = "Page not found"
	ErrPrefixInvalidEmptyInputs     = "The provided inputs must not be empty"
)

// DestinationUnreachable is the error used when the crawler can't reach a destination from its current path.
type DestinationUnreachable struct {
	Destination string
}

// Error returns the string representation of the DestinationUnreachable error.
func (e DestinationUnreachable) Error() string {
	return fmt.Sprintf("%s: %s", ErrPrefixDestinationUnreachable, e.Destination)
}

// PageNotFound is the error used when a request page can't be found in the wiki.
type PageNotFound struct {
	wiki.Page
}

// Error returns the string representation of the PageNotFound error.
func (e PageNotFound) Error() string {
	return fmt.Sprintf("%s: %s", ErrPrefixPageNotFound, e.Title)
}

// InvalidEmptyInput is the error used when the provided inputs are invalid.
type InvalidEmptyInput struct {
	Origin      string
	Destination string
}

// Error returns the string representation of the InvalidEmptyInput error.
func (e InvalidEmptyInput) Error() string {
	return fmt.Sprintf("%s: (%s, %s)", ErrPrefixInvalidEmptyInputs, e.Origin, e.Destination)
}
