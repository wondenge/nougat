package nougat

import (
	"io"
	"net/http"
	"net/url"
)

// Requests

// Request returns a new http.Request created with the Nougat properties.
// Returns any errors parsing the rawURL, encoding query structs, encoding
// the body, or creating the http.Request.
func (r *Nougat) Request() (*http.Request, error) {
	reqURL, err := url.Parse(r.rawURL)
	if err != nil {
		return nil, err
	}

	err = addQueryStructs(reqURL, r.queryStructs)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if r.bodyProvider != nil {
		body, err = r.bodyProvider.Body()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(r.method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	addHeaders(req, r.header)
	return req, err
}
