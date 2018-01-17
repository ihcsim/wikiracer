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

	v sync.Map
}

// NewForward returns an new instance of the Forward crawler.
func NewForward(w wiki.Wiki) *Forward {
	return &Forward{
		Wiki:   w,
		path:   make(chan *wiki.Path),
		errors: make(chan error),
		done:   make(chan struct{}),
		wg:     &sync.WaitGroup{},
		v:      sync.Map{},
	}
}

// Run provides the implementation of the crawling algorithm.
// It doesn't return any values. The result path can be obtained using the Path() method. All errors encountered can be retrieved using the Error() method.
// When all work is completed, the Done() method can be used to signal the caller.
func (f *Forward) Run(origin, destination string) {
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		f.discover(origin, destination, nil)
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

// Done returns a channel which can be used to signify that the crawl is completed.
func (f *Forward) Done() <-chan struct{} {
	return f.done
}

func (f *Forward) discover(origin, destination string, intermediate *wiki.Path) {
	// discover crawls from origin to destination using all the links found in the pages.
	// for every page that it encounters:
	// 1. it marks the page as visited by adding it to the 'v' map
	// 2. it appends the page to the sequence of pages in the 'intermediate' path
	// 3. if the page is the destination page, it returns the 'intermediate' path
	// 4. if the page isn't the destination page and has no links, it returns the 'destination unreachable' error because it has reached a dead end on this path
	// 5. otherwise, it creates a goroutine to crawl every link of the page

	page, err := f.FindPage(origin)
	if err != nil {
		f.errors <- err
		return
	}

	if intermediate == nil {
		intermediate = wiki.NewPath()
	}
	intermediate.AddPage(page)

	if f.visited(origin) {
		f.errors <- errors.LoopDetected{Path: intermediate}
		return
	}
	f.addVisited(origin)

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
			f.discover(link, destination, newPath)
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
