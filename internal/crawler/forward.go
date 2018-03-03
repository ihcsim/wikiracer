package crawler

import (
	"context"
	"sync"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
	"github.com/ihcsim/wikiracer/log"
)

const (
	separator               = "|"
	wikipediaMaxTitlesCount = 50
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
// For every page P that it encounters:
// 1. `P` is appended to the sequence of pages in the _intermediate_ path.
// 2. `P` is marked as a visited page.
// 3. if `P` is the destination page, the _intermediate_ path is returned.
// 4. if `P` isn't the destination page and has no links, the goroutine terminates.
// 5. otherwise, for every link of `P`, the goroutine creates a new goroutine to crawl that linked page.
func (f *Forward) discover(ctx context.Context, titles, destination string, ancestors *wiki.Path) {
	if ctx.Err() != nil {
		log.Instance().Debugf("Canceling crawl operation. Reason=%q", ctx.Err().Error())
		return
	}

	pages, err := f.FindPages(titles, "")
	if err != nil {
		if pageErr, ok := err.(errors.PageNotFound); ok {
			if pageErr.Title != destination {
				return
			}
		}

		log.Instance().Errorf("%s", err)
		f.errors <- err
		return
	}

	for _, page := range pages {
		clonedAncestors := wiki.NewPath()
		if ancestors != nil {
			clonedAncestors.Clone(ancestors)
		}
		clonedAncestors.AddPage(page)

		// skip this page if is previously visited
		if f.visited(page.Title) {
			log.Instance().Debugf("Loop detected. Title=%q Predecessors=%q", page.Title, clonedAncestors)
			continue
		}
		f.addVisited(page.Title)
		log.Instance().Debugf("Found page. Title=%q Predecessors=%q", page.Title, clonedAncestors)

		// found destination
		if page.Title == destination {
			log.Instance().Infof("Found destination. Title=%q Predecessors=%q", page.Title, clonedAncestors)
			f.path <- clonedAncestors
			return
		}

		// this page is a dead end and the racer can't reach the destination from this path.
		if len(page.Links) == 0 {
			log.Instance().Debugf("Dead end page. Title=%q Predecessors=%q", page.Title, clonedAncestors)
			continue
		}

		// Since the Wikipedia API only supports 50 titles in one query,
		// we have to break up the query into multiple calls.
		batchCount := len(page.Links) / wikipediaMaxTitlesCount
		links := make([]string, batchCount+1)
		for index, link := range page.Links {
			// if one of the linked pages is the destination, return it
			if link == destination {
				log.Instance().Infof("Found destination. Title=%q Predecessors=%q", link, clonedAncestors)
				f.addVisited(link)
				clonedAncestors.AddPage(&wiki.Page{Title: link})
				f.path <- clonedAncestors
				return
			}

			links[index/wikipediaMaxTitlesCount] += separator + link
		}

		go func() {
			log.Instance().Debugf("Starting crawl operation. Titles=%q", links)
			for _, link := range links {
				if link == "" {
					continue
				}
				f.discover(ctx, link[1:], destination, clonedAncestors)
			}
		}()
	}
}

func (f *Forward) addVisited(title string) {
	f.v.Store(title, struct{}{})
}

func (f *Forward) visited(title string) bool {
	_, exist := f.v.Load(title)
	return exist
}
