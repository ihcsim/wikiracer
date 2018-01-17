package wikiracer

import (
	"fmt"

	"github.com/ihcsim/wikiracer/errors"
)

// Racer traverses from a wiki page to another using only links.
// It times the traversal journey.
type Racer struct {
	Crawler
	Validator
}

// FindPath attempts to find a path from the origin page to the destination page by traversing all the links that are encountered along the way.
func (r *Racer) FindPath(origin, destination string) string {
	if err := r.Validate(origin, destination); err != nil {
		return fmt.Sprintf("%s", err)
	}

	if origin == destination {
		return origin
	}

	go r.Run(origin, destination)

	for {
		select {
		case path := <-r.Path():
			return fmt.Sprintf("%s", path)
		case err := <-r.Error():
			if _, ok := err.(errors.DestinationUnreachable); !ok {
				return fmt.Sprintf("%s", err)
			}
		case <-r.Done():
			// if we received from done, it means all goroutines completed but no path was found.
			return fmt.Sprintf("%s", errors.DestinationUnreachable{Destination: destination})
		}
	}
}
