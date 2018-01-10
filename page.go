package wikiracer

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
