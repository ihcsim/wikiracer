package wikipedia

// Response is the raw JSON response from the Wikipedia.
type Response struct {
	// Next, if presents, points to the next batch of result.
	Next *NextBatch `json:"continue",omitempty`

	// Result contains pages data received from the Wikipedia. If an URL redirect was performed by wikipedia before the result is retrieved, a 'Redirect' block will be included.
	Result *Query `json:"query",omitempty`

	// Batchcomplete is true if there are no more subsequent batches.
	// Otherwise, it's omitted.
	Batchcomplete bool `json:"",omitempty`

	// Errors is a list of errors returned by the Wikipedia.
	Errors []*ResponseError `json:"",omitempty`

	// Warnings is a list of warnings returned by the Wikipedia.
	Warnings *ResponseWarnings `json:"",omitempty`
}

// NextBatch points to the next batch of result.
// For more information on how 'continue' works, refer to https://www.mediawiki.org/wiki/API:Raw_query_continue
type NextBatch struct {
	// Plcontinue is the title of the first page of the next batch of result.
	Plcontinue string

	// Continue
	Continue string
}

// Query contains pages data received from the Wikipedia.
type Query struct {
	// Redirects represents any URL redirect that Wikipedia performed before the result is retrieved. Wikipedia performs URL redirects for certain pages that may be known by multiple titles.
	// For more information on how 'redirect' works, refer to https://en.wikipedia.org/wiki/Help:Redirect
	Redirects []*Redirect `json:"",omitempty`

	// Pages is the batch of pages received from the Wikipedia.
	Pages []*Page
}

// Redirect represents a single URL redirect performed by Wikipedia. Wikipedia performs URL redirects for certain pages that may be known by multiple titles.
type Redirect struct {
	// From is the user-provided title of the page used in the query.
	From string

	// To is the original title of the page as known by the Wikipedia.
	To string
}

// Page represents a single page as returned by the Wikipedia.
type Page struct {
	// Pageid is the page ID.
	Pageid int

	// Ns is the namespace where the page belongs.
	Ns int

	// Title is the page title.
	Title string

	// Links is the collection of links found in the page.
	Links []Link

	// Missing is true if there is no page with the given title.
	Missing bool `json:'',omitempty`
}

// Link is a link to another page.
type Link struct {
	// Ns is the namespace of the linked page.
	Ns int

	// Title is the title of the linked page.
	Title string
}

// ResponseError is a error returned by the Wikipedia.
type ResponseError struct {
	// Code is the error code.
	Code string

	// Text provides the error message.
	Text string
}

// ResponseWarnings are warnings returned by the Wikipedia.
type ResponseWarnings struct {
	// Main is the main warning.
	Main *MainWarning `json:"",omitempty`

	// Query is query-related warning
	Query *QueryWarning `json:"",omitempty`
}

// MainWarning is the warning returned by the Wikipedia 'main' module.
type MainWarning struct {
	// Warnings are warning messages.
	Warnings string
}

// QueryWarning is the query-related warning
type QueryWarning struct {
	// Warnings are warning messages.
	Warnings string
}
