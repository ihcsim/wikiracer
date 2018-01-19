package wikiracer

import (
	"context"

	"github.com/ihcsim/wikiracer/errors"
)

// Racer traverses from a wiki page to another using only links.
// It times the traversal journey.
type Racer struct {
	Crawler
	Validator
}

// FindPath attempts to find a path from the origin page to the destination page by traversing all the links that are encountered along the way.
// If found, it returns the path from origin to destination.
// Otherwise, if a path isn't found, a DestinationUnreachable error is returned.
// The destination page is considered unreachable if racer can't find it before the context timed out.
// Use ctx to impose timeout on FindPath.
func (r *Racer) FindPath(ctx context.Context, origin, destination string) string {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := r.Validate(origin, destination); err != nil {
		return err.Error()
	}

	if origin == destination {
		return origin
	}

	go r.Run(cancelCtx, origin, destination)

	for {
		select {
		case path := <-r.Path():
			return path.String()

		case err := <-r.Error():
			return err.Error()

		case <-cancelCtx.Done():
			return errors.DestinationUnreachable{Destination: destination}.Error()
		}
	}
}
