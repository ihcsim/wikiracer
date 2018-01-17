package errors

import (
	"fmt"

	"github.com/ihcsim/wikiracer/internal/wiki"
)

// DestinationUnreachable is the error used when the crawler can't reach a destination from its current path.
type DestinationUnreachable struct {
	Destination string
}

// Error returns the string representation of the DestinationUnreachable error.
func (e DestinationUnreachable) Error() string {
	return fmt.Sprintf("%s: %s", "Destination unreachable", e.Destination)
}

// LoopDetected is the error used when the crawler encounters a sequence of pages that form a loop.
type LoopDetected struct {
	Path *wiki.Path
}

// Error is the string representation of the LoopDetected error.
func (e LoopDetected) Error() string {
	return fmt.Sprintf("%s", e.Path)
}

// PageNotFound is the error used when a request page can't be found in the wiki.
type PageNotFound struct {
	wiki.Page
}

// Error returns the string representation of the PageNotFound error.
func (e PageNotFound) Error() string {
	return fmt.Sprintf("%s: %s", "Page not found", e.Title)
}

// InvalidEmptyInput is the error used when the provided inputs are invalid.
type InvalidEmptyInput struct {
	Origin      string
	Destination string
}

// Error returns the string representation of the InvalidEmptyInput error.
func (e InvalidEmptyInput) Error() string {
	return fmt.Sprintf("%s: (%s, %s)", "The provided inputs must not be empty", e.Origin, e.Destination)
}
