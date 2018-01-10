package wikiracer

const (
	requestEndpoint = "https://en.wikipedia.org/w/api.php"
	requestAction   = "query"
	requestProp     = "links"

	pageNamespace = 0

	responseDefaultFormat = "json"
	responseDefaultLimit  = 1
)

var wiki Wiki

// FindPath attempts to find a path from the 'origin' page to the 'destination' page by recursively traversing all the links that are either found in or linked to the 'origin' page.
func FindPath(origin, destination string) string {
	path := make(chan string)
	traverse(origin, destination, path)
	v := <-path
	return origin + " -> " + v
}

func traverse(origin, destination string, path chan string) {
	page, _ := wiki.FindPage(origin)

	for _, link := range page.Links {
		if link == destination {
			path <- link

			// send termination signal to other goroutines?
			return
		}
	}

	for _, link := range page.Links {
		p := make(chan string)
		go traverse(link, destination, p)

		go func() {
			v := <-p
			path <- link + " -> " + v
		}()
	}
}
