package crawler

import (
	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
)

// Forward is a crawler that attempts to find a path from an origin page to a destination page (via links in the pages) by using an uni-directional traversal pattern.
type Forward struct {
	wiki.Wiki
}

// NewForward returns an new instance of the Forward crawler.
func NewForward(w wiki.Wiki) *Forward {
	return &Forward{Wiki: w}
}

// Discover provides the crawling algorithm of the Forward Crawler.
// If found, it returns the path from the origin to the destination.
// Otherwise, if such a path doesn't exist, it returns an error.
func (f *Forward) Discover(origin, destination string) (*wiki.Path, error) {
	page, err := f.FindPage(origin)
	if err != nil {
		return nil, err
	}

	if len(page.Links) == 0 {
		// this page is a dead end and the racer can't reach the destination from this path.
		return nil, errors.DestinationUnreachable{Destination: destination}
	}

	// create a goroutine for each link of this page.
	// every goroutine is assigned its own path and error channels.
	// the path channel receives either the path to the destination or nothing at all if the destination is unreachable on this path.
	// the error channel receives all errors returned by the goroutines, including the 'destination unreachable' error.
	var (
		pathChans = make([]chan *wiki.Path, len(page.Links))
		errChans  = make([]chan error, len(page.Links))
	)
	for i := 0; i < len(page.Links); i++ {
		pathChans[i] = make(chan *wiki.Path)
		errChans[i] = make(chan error)
	}

	for index, link := range page.Links {
		// found destination
		if link == destination {
			dp, err := f.FindPage(link)
			if err != nil {
				return nil, err
			}
			return &wiki.Path{dp}, nil
		}

		go func(link string, index int) {
			// follow a new path using 'link' as the starting point.
			// returned values of all recursive calls are captured here.
			path, err := f.Discover(link, destination)
			if err != nil {
				errChans[index] <- err
				return
			}

			// create a new path by concatenating the returned path to the
			// current 'link' page.
			dp, err := f.FindPage(link)
			if err != nil {
				errChans[index] <- err
			}

			parent := &wiki.Path{dp}
			parent.Concat(path)
			pathChans[index] <- parent
		}(link, index)
	}

	// keep looping until either:
	// i.   a path to the destination is received from one of the goroutines,
	// ii.  an error is received from one of the goroutines, or
	// iii. all goroutines have returned the 'destination unreachable' error
	goroutinesCount := len(page.Links)
	for {
		for i := 0; i < len(page.Links); i++ {
			// the pathChans[i] and errChans[i] are set to nil after we are done with them.
			// we can't close them, because closed channels are always ready to receive, where
			// nil channels are always block. Refer https://dave.cheney.net/2013/04/30/curious-channels
			if pathChans[i] != nil || errChans[i] != nil {
				select {
				case path := <-pathChans[i]:
					pathChans[i] = nil
					return path, nil

				case err := <-errChans[i]:
					errChans[i] = nil
					cast, ok := err.(errors.DestinationUnreachable)
					if !ok {
						return nil, err
					}

					goroutinesCount--
					if goroutinesCount == 0 {
						// all the goroutines have returned the 'destination unreachable' error
						return nil, cast
					}

				default:
				}
			}
		}
	}
}
