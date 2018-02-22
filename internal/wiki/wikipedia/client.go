package wikipedia

import (
	"encoding/json"

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

// FindPage returns the page of the given title.
func (c *Client) FindPage(title, nextBatch string) (*wiki.Page, error) {
	response, err := c.query(title, nextBatch)
	if err != nil {
		return nil, err
	}

	page := response.Result.Pages[0]
	if page.Missing {
		return nil, errors.PageNotFound{wiki.Page{Title: title}}
	}

	var links []string
	for _, link := range page.Links {
		links = append(links, link.Title)
	}
	result := &wiki.Page{
		ID:        response.Result.Pages[0].Pageid,
		Title:     response.Result.Pages[0].Title,
		Namespace: response.Result.Pages[0].Ns,
		Links:     links,
	}

	// the links in a page are usually returned in batches.
	// when `batchcomplete` is set in the response, it implies that the server has returned the last batch of links for this page.
	// when `plcontinue` is set in the response, it implies that there are more links yet to be fetched.

	if response.Batchcomplete {
		return result, nil
	}

	if response.Next != nil && response.Next.Plcontinue != "" {
		partial, err := c.FindPage(title, response.Next.Plcontinue)
		if err != nil {
			return nil, err
		}

		result.Links = append(result.Links, partial.Links...)
	}

	return result, nil
}

func (c *Client) query(title, plcontinue string) (*Response, error) {
	query := map[string]string{
		"action":        "query",
		"prop":          "links",
		"format":        responseFormat,
		"formatversion": responseFormatVersion,
		"pllimit":       responseLimits,
		"plnamespace":   namespace,
		"titles":        title,
		"redirects":     "true",
		"utf8":          "true",
	}

	if plcontinue != "" {
		query["plcontinue"] = plcontinue
	}

	content, err := c.api(c.client, query)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
