package wikiracer

import (
	"fmt"
	"time"
)

// Result captures the duration to discover the path from the origin page to the destination page.
type Result struct {
	// Path represents an ordered sequence of pages from the origin page to the destination page.
	Path []byte

	// Duration captures the time taken to discover path.
	Duration time.Duration

	// Err are errors captured during the path discovery.
	Err error
}

// String returns a string representation of a result.
func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("%s", r.Err)
	}

	return fmt.Sprintf("Path: %q, Duration: %s", r.Path, r.Duration)
}
