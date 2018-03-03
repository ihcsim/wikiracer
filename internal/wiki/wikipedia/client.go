package wikipedia

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
	"github.com/sadbox/mediawiki"
)

const (
	url                   = "https://en.wikipedia.org/w/api.php"
	userAgent             = "wikiracer"
	responseFormat        = "json"
	responseFormatVersion = "2"
	responseLimits        = "max"
	namespace             = "0"

	wikipediaTooManyRequestsErr = "Error: 429, Too Many Requests"
	coolDownDuration            = time.Second
)

// Client can communicate with the Wikipedia URL.
type Client struct {
	client *mediawiki.MWApi
	api    apiFunc
}

// NewClient creates a new instanc of Client.
func NewClient() (*Client, error) {
	c, err := mediawiki.New(url, userAgent)
	return &Client{
		client: c,
		api:    (*mediawiki.MWApi).API,
	}, err
}

type apiFunc func(api *mediawiki.MWApi, values ...map[string]string) ([]byte, error)

// FindPages returns the page of the given title.
func (c *Client) FindPages(titles, nextBatch string) ([]*wiki.Page, error) {
	response, err := c.query(titles, nextBatch)
	if err != nil {
		return nil, err
	}

	if response.Errors != nil {
		err := handleErrors(response.Errors)
		return nil, err
	}

	if response.Warnings != nil {
		err := handleWarnings(response.Warnings)
		return nil, err
	}

	results := []*wiki.Page{}

	if response.Result != nil {
		for _, page := range response.Result.Pages {
			if page.Missing {
				return nil, errors.PageNotFound{wiki.Page{Title: page.Title}}
			}

			var links []string
			for _, link := range page.Links {
				links = append(links, link.Title)
			}

			results = append(results, &wiki.Page{
				ID:        page.Pageid,
				Title:     page.Title,
				Namespace: page.Ns,
				Links:     links,
			})
		}
	}

	// the links in a page are usually returned in batches.
	// when `batchcomplete` is set in the response, it implies that the server has returned the last batch of links for this page.
	// when `plcontinue` is set in the response, it implies that there are more links yet to be fetched.

	if response.Batchcomplete {
		return results, nil
	}

	if response.Next != nil && response.Next.Plcontinue != "" {
		nextBatch, err := c.FindPages(titles, response.Next.Plcontinue)
		if err != nil {
			return nil, err
		}

		for _, batchResult := range nextBatch {
			for _, result := range results {
				if result.ID == batchResult.ID {
					result.Links = append(result.Links, batchResult.Links...)
				}
			}
		}
	}

	return results, nil
}

func (c *Client) query(titles, plcontinue string) (*Response, error) {
	query := map[string]string{
		"action":        "query",
		"prop":          "links",
		"format":        responseFormat,
		"formatversion": responseFormatVersion,
		"pllimit":       responseLimits,
		"plnamespace":   namespace,
		"titles":        titles,
		"redirects":     "true",
		"utf8":          "true",
	}

	if plcontinue != "" {
		query["plcontinue"] = plcontinue
	}

	var (
		content []byte
		err     error
	)
	for {
		content, err = c.api(c.client, query)
		if err != nil {
			return nil, err
		}

		// check if the wikipedia API returns a 429 error
		if !strings.Contains(string(content), wikipediaTooManyRequestsErr) {
			break
		}

		// retry the API call after the cooldown duration expires
		time.Sleep(coolDownDuration)
	}

	var response Response
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleErrors(errors []*ResponseError) error {
	err := &serverError{}
	for _, e := range errors {
		err.msg += fmt.Sprintf("%s\n", e.Text)
	}

	return err
}

func handleWarnings(warnings *ResponseWarnings) error {
	err := &serverError{}
	if warnings.Main != nil {
		err.msg = warnings.Main.Warnings
	}

	if warnings.Query != nil {
		err.msg += warnings.Query.Warnings
	}

	return err
}
