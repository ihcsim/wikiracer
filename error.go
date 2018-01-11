package wikiracer

import "fmt"

const (
	ErrPrefixPageNotFound       = "Page not found"
	ErrPrefixInvalidEmptyInputs = "The provided inputs must not be empty"
)

// PageNotFound is the error used when a non-existent page is requested.
type PageNotFound struct {
	Page
}

// Error returns the string representation of the PageNotFound error.
func (p PageNotFound) Error() string {
	return fmt.Sprintf("%s: %s", ErrPrefixPageNotFound, p.Title)
}

// InvalidEmptyInput is the error used when the provided inputs are empty.
type InvalidEmptyInput struct {
	origin      string
	destination string
}

// Error returns the string representation of the InvalidEmptyInput error.
func (i InvalidEmptyInput) Error() string {
	return fmt.Sprintf("%s: (%s, %s)", ErrPrefixInvalidEmptyInputs, i.origin, i.destination)
}
