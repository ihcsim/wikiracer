package wikiracer

import (
	"fmt"

	"github.com/ihcsim/wikiracer/internal/wiki"
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

	go func() {
		intermediate := wiki.NewPath()
		r.Discover(origin, destination, intermediate)
	}()

	var errors = []error{}
	for {
		select {
		case path := <-r.Path():
			return fmt.Sprintf("%s", path)
		case err := <-r.Error():
			errors = append(errors, err)
		}
	}
}
