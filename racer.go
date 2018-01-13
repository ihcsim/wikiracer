package wikiracer

import (
	"fmt"
	"sync"
)

var wiki Wiki

// FindPath attempts to find a path from the 'origin' page to the 'destination' page by recursively traversing all the links that are either found in or linked to the 'origin' page.
func FindPath(origin, destination string) string {
	if validationErr := validate(origin, destination); validationErr != nil {
		return fmt.Sprintf("%s", validationErr)
	}

	if origin == destination {
		return origin
	}

	var (
		path = make(chan string)
		err  = make(chan error)
		wg   = &sync.WaitGroup{}
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		p, e := traverse(origin, destination, wg)
		if e != nil {
			err <- e
			return
		}
		path <- p
	}()

	go func() {
		defer func() {
			close(path)
			close(err)
		}()

		wg.Wait()
		//all goroutines exited
	}()

	for {
		select {
		case v := <-path:
			return origin + " -> " + v
		case e := <-err:
			if t, ok := e.(NoLinksFound); ok {
				t.origin = origin
				return fmt.Sprintf("%s", t)
			}

			return fmt.Sprintf("%s", e)
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

func traverse(origin, destination string, wg *sync.WaitGroup) (string, error) {
	page, wikiErr := wiki.FindPage(origin)
	if wikiErr != nil {
		return "", wikiErr
	}

	if len(page.Links) == 0 {
		return "", nil
	}

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
			return link, nil
		}

		wg.Add(1)
		go func(link string, index int) {
			defer wg.Done()

			p, e := traverse(link, destination, wg)
			if e != nil {
				errChans[index] <- e
				return
			}

			if p != "" {
				p = link + " -> " + p
			}

			pathChans[index] <- p
		}(link, index)
	}

	stillAlive := len(page.Links)
	for {
		for i := 0; i < len(page.Links); i++ {
			select {
			case v := <-pathChans[i]:
				if v != "" {
					return v, nil
				}

				stillAlive -= 1
				if stillAlive == 0 {
					return "", NoLinksFound{destination: destination}
				}

			case e := <-errChans[i]:
				return "", e
			default:
			}
		}
	}
}
