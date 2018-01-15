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

// FindPath attempts to find a path from the 'origin' page to the 'destination' page by traversing all the links that are found either in the 'origin' page or its linked pages.
func (r *Racer) FindPath(origin, destination string) string {
	if err := r.Validate(origin, destination); err != nil {
		return fmt.Sprintf("%s", err)
	}

	if origin == destination {
		return origin
	}

	var (
		pathChan = make(chan *wiki.Path)
		errChan  = make(chan error)
	)

	defer func() {
		close(pathChan)
		close(errChan)
	}()

	// begin the link traversal process in a goroutine
	go func() {
		path, err := r.Discover(origin, destination)
		if err != nil {
			errChan <- err
			return
		}
		pathChan <- path
	}()

	// wait for results to arrive via channels
	for {
		select {
		case path := <-pathChan:
			return origin + " -> " + fmt.Sprintf("%s", path)
		case err := <-errChan:
			return fmt.Sprintf("%s", err)
		}
	}
}
