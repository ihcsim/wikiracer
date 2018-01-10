package wikiracer

// Wiki provides a collection of methods to communicate with a Wiki instance.
type Wiki interface {

	// FindPage returns the page of the given title.
	FindPage(title string) (*Page, error)
}
