package wikiracer

import (
	"github.com/ihcsim/wikiracer/internal/wiki"
)

// Crawler can find a path from one wiki page to another.
// The starting page is the origin. The target page is the destination.
type Crawler interface {
	// Run provides the implementation of the crawling algorithm.
	// It doesn't return any values. The result path can be obtained using the Path() method. All errors encountered can be retrieved using the Error() method.
	// When all work is completed, the Done() method can be used to signal the caller.
	Run(origin, destination string)

	// Path returns a channel which can be used to receive the path result.
	Path() <-chan *wiki.Path

	// Error returns a channel which can be used to receive any errors encountered by Discover().
	Error() <-chan error

	// Done returns a channel which can be used to signify that the crawl is completed.
	Done() <-chan struct{}
}
