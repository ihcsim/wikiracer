package wikiracer

import "fmt"

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
	if validationErr := validate(origin, destination); validationErr != nil {
		return fmt.Sprintf("%s", validationErr)
	}

	var (
		fullPath = make(chan string)
		err      = make(chan error)
		result   = ""
	)

	go traverse(origin, destination, fullPath, err)

	select {
	case result = <-fullPath:
		return origin + " -> " + result
	case e := <-err:
		return fmt.Sprintf("%s", e)
	}
}

func validate(origin, destination string) error {
	if len(origin) == 0 || len(destination) == 0 {
		return InvalidEmptyInput{origin: origin, destination: destination}
	}

	if _, err := wiki.FindPage(origin); err != nil {
		return err
	}

	if _, err := wiki.FindPage(destination); err != nil {
		return err
	}

	return nil
}

func traverse(origin, destination string, path chan string, err chan error) {
	if origin == destination {
		path <- origin
		return
	}

	page, wikiErr := wiki.FindPage(origin)
	if wikiErr != nil {
		err <- wikiErr
		return
	}

	if len(page.Links) == 0 {
		return
	}

	for _, link := range page.Links {
		if link == destination {
			path <- link
			return
		}
	}

	for _, link := range page.Links {
		intermediate := make(chan string)
		go traverse(link, destination, intermediate, err)

		go func() {
			v := <-intermediate
			path <- link + " -> " + v
		}()
	}
}
