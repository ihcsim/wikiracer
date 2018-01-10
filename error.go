package wikiracer

import "fmt"

const pageNotFoundErrorPrefix = "Page not found"

// PageNotFound is the error used when a non-existent page is requested.
type PageNotFound struct {
	*Page
}

// Error returns the string representation of the PageNotFound error.
func (p PageNotFound) Error() string {
	return fmt.Sprintf("%s: %s", pageNotFoundErrorPrefix, p.Title)
}
