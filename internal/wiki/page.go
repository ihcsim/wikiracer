package wiki

import (
	"strings"
	"sync"
)

const pathDelimiter = " -> "

// Page represents a wiki page
type Page struct {
	// ID is the page's ID.
	ID int

	// Title is the page's title.
	Title string

	// Namespace is the page's namespace.
	Namespace int

	// Links is the collection of all the links (to other pages) found in the page.
	Links []string
}

// Path is an ordered sequence of pages which forms a path from the first page to the last page.
type Path struct {
	mux      sync.Mutex
	sequence []*Page
}

// NewPath returns a new instance of path.
func NewPath() *Path {
	return &Path{
		mux:      sync.Mutex{},
		sequence: []*Page{},
	}
}

// AddPage appends a page to the path.
func (p *Path) AddPage(page *Page) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.sequence = append(p.sequence, page)
}

// String returns the string representation of the path.
func (p *Path) String() string {
	p.mux.Lock()
	defer p.mux.Unlock()

	var s string
	for _, page := range p.sequence {
		s += page.Title + pathDelimiter
	}

	i := strings.LastIndex(s, pathDelimiter)
	return s[:i]
}

// Clone copies the sequence of p2 into p1.
func (p *Path) Clone(p2 *Path) int {
	p.sequence = make([]*Page, len(p2.sequence))
	return copy(p.sequence, p2.sequence)
}
