package wikiracer

import (
	"context"
	"strings"
	"time"

	"github.com/ihcsim/wikiracer/errors"
)

// WikiRacer traverses from a wiki page to another using only links.
// It times the traversal journey.
type WikiRacer struct {
	Crawler
	Validator
}

// New returns a new instance of the WikiRacer.
func New(c Crawler, v Validator) *WikiRacer {
	return &WikiRacer{Crawler: c, Validator: v}
}

// TimedFindPath captues the duration to crawl from the origin page to the destination page.
// If a timeout is specified in ctx, then the result duration will not exceed the timeout.
func (r *WikiRacer) TimedFindPath(ctx context.Context, origin, destination string) *Result {
	start := time.Now()
	result := r.FindPath(ctx, origin, destination)
	end := time.Now()

	result.Duration = end.Sub(start)
	return result
}

// FindPath attempts to find a path from the origin page to the destination page by traversing all the links that are encountered along the way.
// If found, it returns the path from origin to destination.
// Otherwise, if a path isn't found, a DestinationUnreachable error is returned.
// The destination page is considered unreachable if racer can't find it before the context timed out.
// Use ctx to impose timeout on FindPath.
func (r *WikiRacer) FindPath(ctx context.Context, origin, destination string) *Result {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := r.Validate(origin, destination); err != nil {
		return &Result{Err: err}
	}

	if origin == destination {
		return &Result{Path: strings.NewReader(origin)}
	}

	go r.Run(cancelCtx, origin, destination)

	for {
		select {
		case path := <-r.Path():
			return &Result{Path: strings.NewReader(path.String())}

		case err := <-r.Error():
			return &Result{Err: err}

		case <-cancelCtx.Done():
			return &Result{Err: errors.DestinationUnreachable{Destination: destination}}
		}
	}
}
