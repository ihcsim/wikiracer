package wiki

// Wiki provides a collection of methods to communicate with a wiki instance.
type Wiki interface {

	// FindPages returns the page of the given title.
	FindPages(titles, nextBatch string) ([]*Page, error)
}
