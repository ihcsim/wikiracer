package crawler

import (
	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
)

// Forward is a crawler that attempts to find a path from an origin page to a destination page using an uni-directional traversal pattern.
type Forward struct {
	wiki.Wiki
	path   chan *wiki.Path
	errors chan error
}

// NewForward returns an new instance of the Forward crawler.
func NewForward(w wiki.Wiki) *Forward {
	return &Forward{
		Wiki:   w,
		path:   make(chan *wiki.Path),
		errors: make(chan error),
	}
}

// Path returns a channel which receives the path result from the children goroutines.
func (f *Forward) Path() <-chan *wiki.Path {
	return f.path
}

// Error returns a channel which receives errors from the children goroutines.
func (f *Forward) Error() <-chan error {
	return f.errors
}

// Discover provides the crawling implementation of a Crawler.
// The intermediate path struct is used to keep track of all the pages encountered so far.
// If found, the path from the origin page to the destination page can be retrieved using the Path() method.
// Otherwise, if such a path doesn't exist, the Error() method will return a DestinationUnreachable error.
func (f *Forward) Discover(origin, destination string, intermediate *wiki.Path) {
	page, err := f.FindPage(origin)
	if err != nil {
		f.errors <- err
		return
	}

	intermediate.AddPage(page)

	// found destination
	if page.Title == destination {
		f.path <- intermediate
		return
	}

	// this page is a dead end and the racer can't reach the destination from this path.
	if len(page.Links) == 0 {
		f.errors <- errors.DestinationUnreachable{Destination: destination}
		return
	}

	for _, link := range page.Links {
		go func(link string) {
			newPath := wiki.NewPath()
			newPath.Clone(intermediate)
			f.Discover(link, destination, newPath)
		}(link)
	}
}

// Exit closes the path and error channels.
func (f *Forward) Exit() {
}
