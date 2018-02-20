package wikipedia

import (
	"encoding/json"
	"fmt"

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
	*mediawiki.MWApi
}

// New creates a new instanc of Client.
func New() (*Client, error) {
	api, err := mediawiki.New(url, userAgent)
	return &Client{api}, err
}

// FindPage returns the page of the given title.
func (c *Client) FindPage(title string) (*wiki.Page, error) {
	response, err := c.api(title)
	if err != nil {
		return nil, err
	}

	for _, page := range response.Result.Pages {
		for _, link := range page.Links {
			fmt.Printf("%+v\n", link)
		}
	}

	// handle continue

	// handle batchcomplete

	return nil, nil
}

func (c *Client) api(title string) (*Response, error) {
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
	content, err := c.API(query)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
