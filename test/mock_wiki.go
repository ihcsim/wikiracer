package test

import (
	"strings"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
)

const separator = "|"

// MockWiki is an in-memory wiki
type MockWiki struct {
	pages map[string]*wiki.Page
}

// NewMockWiki returns a new instance of MockWiki
func NewMockWiki() *MockWiki {
	testData := map[string]*wiki.Page{
		"1984 Summer Olympics": &wiki.Page{ID: 2000, Title: "1984 Summer Olympics", Namespace: 0, Links: []string{"7-Eleven", "Afghanistan"}},
		"2010 Winter Olympics": &wiki.Page{ID: 2009, Title: "2010 Winter Olympics", Namespace: 0, Links: []string{"1984 Summer Olympics"}},
		"7-Eleven":             &wiki.Page{ID: 2001, Title: "7-Eleven", Namespace: 0, Links: []string{"Big C", "Calgary", "Eurocash"}},
		"Afghanistan":          &wiki.Page{ID: 2002, Title: "Afghanistan", Namespace: 0, Links: []string{}},
		"Alexander the Great":  &wiki.Page{ID: 1000, Title: "Alexander the Great", Namespace: 0, Links: []string{"Apepi", "Greek language", "Diodotus I"}},
		"Apepi":                &wiki.Page{ID: 1005, Title: "Apepi", Namespace: 0},
		"Big C":                &wiki.Page{ID: 2003, Title: "Big C", Namespace: 0, Links: []string{"Vancouver"}},
		"Calgary":              &wiki.Page{ID: 2004, Title: "Calgary", Namespace: 0},
		"Eurocash":             &wiki.Page{ID: 2005, Title: "Eurocash", Namespace: 0, Links: []string{"Małpka Express", "Tea"}},
		"Diodotus I":           &wiki.Page{ID: 1007, Title: "Diodotus I", Namespace: 0},
		"Fruit anatomy":        &wiki.Page{ID: 1001, Title: "Fruit anatomy", Namespace: 0, Links: []string{"Segment"}},
		"Greek language":       &wiki.Page{ID: 1002, Title: "Greek language", Namespace: 0, Links: []string{"Fruit anatomy"}},
		"Małpka Express":       &wiki.Page{ID: 2006, Title: "Małpka Express", Namespace: 0},
		"Michael Jordan":       &wiki.Page{ID: 1006, Title: "Michael Jordan", Namespace: 0},
		"Mike Tyson":           &wiki.Page{ID: 1003, Title: "Mike Tyson", Namespace: 0, Links: []string{"Alexander the Great", "1984 Summer Olympics"}},
		"Segment":              &wiki.Page{ID: 1004, Title: "Segment", Namespace: 0, Links: []string{"Vancouver"}},
		"Tea":                  &wiki.Page{ID: 2007, Title: "Tea", Namespace: 0},
		"Vancouver":            &wiki.Page{ID: 2008, Title: "Vancouver", Namespace: 0, Links: []string{"2010 Winter Olympics"}},
	}
	return &MockWiki{pages: testData}
}

// FindPages returns the page with the given title, if it exists.
// Otherwise, it returns a 'page not found' error.
func (m *MockWiki) FindPages(titles, nextBatch string) ([]*wiki.Page, error) {
	pages := []*wiki.Page{}
	for _, title := range strings.Split(titles, separator) {
		page, exist := m.pages[title]
		if !exist {
			return nil, errors.PageNotFound{wiki.Page{Title: title}}
		}

		pages = append(pages, page)
	}

	return pages, nil
}
