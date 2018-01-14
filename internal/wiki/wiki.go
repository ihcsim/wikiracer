package wiki

// Wiki provides a collection of methods to communicate with a wiki instance.
type Wiki interface {

	// FindPage returns the page of the given title.
	FindPage(title string) (*Page, error)
}
