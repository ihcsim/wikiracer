package crawler

import (
	"context"
	"sync"

	"github.com/ihcsim/wikiracer/internal/wiki"
	"github.com/ihcsim/wikiracer/log"
)

// Forward is a crawler that attempts to find a path from an origin page to a destination page using an uni-directional traversal pattern.
type Forward struct {
	wiki.Wiki
	path   chan *wiki.Path
	errors chan error
	v      sync.Map
}

// NewForward returns an new instance of the Forward crawler.
func NewForward(w wiki.Wiki) *Forward {
	return &Forward{
		Wiki:   w,
		path:   make(chan *wiki.Path),
		errors: make(chan error),
		v:      sync.Map{},
	}
}

// Run provides the implementation of the crawling algorithm.
// The result path can be obtained using the Path() method.
// All errors encountered can be retrieved using the Error() method.
// ctx can be used to impose timeout on Run.
func (f *Forward) Run(ctx context.Context, origin, destination string) {
	go f.discover(ctx, origin, destination, nil)
}

// Path returns a channel which receives the path result from the children goroutines.
func (f *Forward) Path() <-chan *wiki.Path {
	return f.path
}

// Error returns a channel which receives errors from the children goroutines.
func (f *Forward) Error() <-chan error {
	return f.errors
}

// discover crawls from origin to destination using all the links found in the pages.
// For every page, P, that it encounters:
// 1. it appends the P to the sequence of pages in the 'intermediate' path
// 2. it marks the P as visited by adding it to the 'v' map
// 3. if P is the destination page, it returns the 'intermediate' path
// 4. if P isn't the destination page and has no links, it's ignored. The goroutine is terminated.
// 5. otherwise, for every link of P, the goroutine creates a new goroutine to crawl that linked page.
func (f *Forward) discover(ctx context.Context, origin, destination string, intermediate *wiki.Path) {
	if ctx.Err() != nil {
		log.Instance().Debugf("Canceling crawl operation. Title=%q Reason=%q", origin, ctx.Err().Error())
		return
	}

	page, err := f.FindPage(origin)
	if err != nil {
		log.Instance().Errorf("%s", err)
		f.errors <- err
		return
	}

	if intermediate == nil {
		intermediate = wiki.NewPath()
	}
	intermediate.AddPage(page)

	if f.visited(origin) {
		log.Instance().Warningf("Loop detected. Title=%q Predecessors=%q", page.Title, intermediate)
		return
	}
	f.addVisited(origin)
	log.Instance().Infof("Found page. Title=%q Predecessors=%q", page.Title, intermediate)

	// found destination
	if page.Title == destination {
		log.Instance().Infof("Found destination. Title=%q Predecessors=%q", page.Title, intermediate)
		f.path <- intermediate
		return
	}

	// this page is a dead end and the racer can't reach the destination from this path.
	if len(page.Links) == 0 {
		log.Instance().Noticef("Dead end page. Title=%q Predecessors=%q", page.Title, intermediate)
		return
	}

	for _, link := range page.Links {
		go func(link string) {
			log.Instance().Debugf("Starting crawl operation. Title=%q", link)
			go func() {
				newPath := wiki.NewPath()
				newPath.Clone(intermediate)
				f.discover(ctx, link, destination, newPath)
			}()

			select {
			case <-ctx.Done():
				log.Instance().Debugf("Finishing crawl operation. Title=%q Reason=%q", link, ctx.Err().Error())
			}
		}(link)
	}
}

func (f *Forward) addVisited(title string) {
	f.v.Store(title, struct{}{})
}

func (f *Forward) visited(title string) bool {
	_, exist := f.v.Load(title)
	return exist
}
