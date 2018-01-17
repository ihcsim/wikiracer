package crawler

import (
	"sync"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
)

// Forward is a crawler that attempts to find a path from an origin page to a destination page using an uni-directional traversal pattern.
type Forward struct {
	wiki.Wiki
	path   chan *wiki.Path
	errors chan error
	done   chan struct{}
	wg     *sync.WaitGroup
}

// NewForward returns an new instance of the Forward crawler.
func NewForward(w wiki.Wiki) *Forward {
	return &Forward{
		Wiki:   w,
		path:   make(chan *wiki.Path),
		errors: make(chan error),
		done:   make(chan struct{}),
		wg:     &sync.WaitGroup{},
	}
}

// Run provides the implementation of the crawling algorithm.
// It doesn't return any values. The result path can be obtained using the Path() method. All errors encountered can be retrieved using the Error() method.
// When all work is completed, the Done() method can be used to signal the caller.
func (f *Forward) Run(origin, destination string) {
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		f.Discover(origin, destination, nil)
	}()

	f.wg.Wait()
	f.done <- struct{}{}
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
// The intermediate path struct is used to keep track of all the pages encountered so far. It can be set to nil for the first call to Discover.
// If found, the path from the origin page to the destination page can be retrieved using the Path() method.
// Otherwise, if such a path doesn't exist, the Error() method will return a DestinationUnreachable error.
func (f *Forward) Discover(origin, destination string, intermediate *wiki.Path) {
	page, err := f.FindPage(origin)
	if err != nil {
		f.errors <- err
		return
	}

	if intermediate == nil {
		intermediate = wiki.NewPath()
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
		f.wg.Add(1)
		go func(link string) {
			defer f.wg.Done()
			newPath := wiki.NewPath()
			newPath.Clone(intermediate)
			f.Discover(link, destination, newPath)
		}(link)
	}
}

// Done returns a channel which can be used to signify that the crawl is completed.
func (f *Forward) Done() <-chan struct{} {
	return f.done
}
