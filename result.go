package wikiracer

import (
	"io"
	"time"
)

// Result captures the duration to discover the path from the origin page to the destination page.
type Result struct {
	// Path represents an ordered sequence of pages from the origin page to the destination page.
	Path io.Reader

	// Duration captures the time taken to discover path.
	Duration time.Duration

	// Err are errors captured during the path discovery.
	Err error
}
