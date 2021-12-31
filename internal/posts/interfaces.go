package posts

import "net/http"

// HTTPClient is interface for executing http requests
// http.DefaultClient implements it
type HTTPClient interface {
	// Do executes http request and retruns response
	Do(*http.Request) (*http.Response, error)
}
