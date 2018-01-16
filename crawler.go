package wikiracer

import "github.com/ihcsim/wikiracer/internal/wiki"

// Crawler can find a path from one wiki page to another.
// The starting page is the origin. The target page is the destination.
type Crawler interface {

	// Discover provides the crawling implementation of a Crawler.
	// The intermediate path struct is used to keep track of all the pages encountered so far.
	// If found, the path from the origin page to the destination page can be retrieved using the Path() method.
	// Otherwise, if such a path doesn't exist, the Error() method will return a DestinationUnreachable error.
	Discover(origin, destination string, intermediate *wiki.Path)

	// Path returns a channel which receives the path result from the children goroutines.
	Path() <-chan *wiki.Path

	// Error returns a channel which receives errors from the children goroutines.
	Error() <-chan error

	// Exit closes the path and error channels.
	Exit()
}
