package wikiracer

import "fmt"

var wiki Wiki

// FindPath attempts to find a path from the 'origin' page to the 'destination' page by recursively traversing all the links that are found either in the 'origin' page or its linked pages.
func FindPath(origin, destination string) string {
	if err := validate(origin, destination); err != nil {
		return fmt.Sprintf("%s", err)
	}

	if origin == destination {
		return origin
	}

	var (
		pathChan = make(chan string)
		errChan  = make(chan error)
	)

	defer func() {
		close(pathChan)
		close(errChan)
	}()

	// begin the link traversal process in a goroutine
	go func() {
		path, err := traverse(origin, destination)
		if err != nil {
			errChan <- err
			return
		}
		pathChan <- path
	}()

	// wait for results to arrive via channels
	for {
		select {
		case path := <-pathChan:
			return origin + " -> " + path
		case err := <-errChan:
			if t, ok := err.(DestinationUnreachable); ok {
				return fmt.Sprintf("%s", t)
			}

			return fmt.Sprintf("%s", err)
		}
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

func traverse(origin, destination string) (string, error) {
	page, err := wiki.FindPage(origin)
	if err != nil {
		return "", err
	}

	if len(page.Links) == 0 {
		// this page is a dead end and the racer can't reach the destination from this path.
		return "", DestinationUnreachable{destination: destination}
	}

	// create a goroutine for each link of this page.
	// every goroutine is assigned its own path and error channels.
	// the path channel receives either the path to the destination or nothing at all if the destination is unreachable on this path.
	// the error channel receives all errors returned by the goroutines, including the 'destination unreachable' error.
	var (
		pathChans = make([]chan string, len(page.Links))
		errChans  = make([]chan error, len(page.Links))
	)
	for i := 0; i < len(page.Links); i++ {
		pathChans[i] = make(chan string)
		errChans[i] = make(chan error)
	}

	for index, link := range page.Links {
		if link == destination {
			// found destination
			return link, nil
		}

		go func(link string, index int) {
			// follow a new path using 'link' as the starting point.
			// returned values of all recursive calls are captured here.
			path, err := traverse(link, destination)
			if err != nil {
				errChans[index] <- err
				return
			}

			pathChans[index] <- link + " -> " + path
		}(link, index)
	}

	// keep looping until either:
	// i.   a path to the destination is received from one of the goroutines,
	// ii.  an error is received from one of the goroutines, or
	// iii. all goroutines have returned the 'destination unreachable' error
	goroutinesCount := len(page.Links)
	for {
		for i := 0; i < len(page.Links); i++ {
			// the pathChans[i] and errChans[i] are set to nil after we are done with them.
			// we can't close them, because closed channels are always ready to receive, where
			// nil channels are always block. Refer https://dave.cheney.net/2013/04/30/curious-channels
			if pathChans[i] != nil || errChans[i] != nil {
				select {
				case path := <-pathChans[i]:
					pathChans[i] = nil
					return path, nil

				case err := <-errChans[i]:
					errChans[i] = nil
					cast, ok := err.(DestinationUnreachable)
					if !ok {
						return "", err
					}

					goroutinesCount -= 1
					if goroutinesCount == 0 {
						// all the goroutines have returned the 'destination unreachable' error
						return "", cast
					}

				default:
				}
			}
		}
	}
}
