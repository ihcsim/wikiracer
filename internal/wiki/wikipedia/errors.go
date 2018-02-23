package wikipedia

// serverError represents an 'invalid API action' error.
type serverError struct {
	msg string
}

// Error returns the string representation of the error.
func (e *serverError) Error() string {
	return e.msg
}
