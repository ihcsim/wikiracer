package wiki

import "strings"

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

// Path is an ordered sequence of pages where every page (except for the last page) is linked to the next page via its wiki link.
type Path []*Page

// Concat joins a given path to this page.
func (p *Path) Concat(path *Path) {
	for _, page := range *path {
		p.AddPage(page)
	}
}

// AddPage appends a page to the path.
func (p *Path) AddPage(page *Page) {
	(*p) = append((*p), page)
}

func (p *Path) String() string {
	var s string
	for _, page := range *p {
		s += page.Title + pathDelimiter
	}

	i := strings.LastIndex(s, pathDelimiter)
	return s[:i]
}
