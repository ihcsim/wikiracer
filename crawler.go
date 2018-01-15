package wikiracer

import "github.com/ihcsim/wikiracer/internal/wiki"

// Crawler can find a path from one wiki page to another.
// The starting page is the origin. The target page is the destination.
type Crawler interface {

	// Discover provides the crawling implementation of a Crawler.
	// If found, it returns the path from origin to destination.
	// Otherwise, if such a path doesn't exist, it returns an empty string and an error.
	Discover(origin, destination string) (*wiki.Path, error)
}
