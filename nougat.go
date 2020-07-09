package nougat

import (
	"net/http"
)

// Doer executes http requests.  It is implemented by *http.Client.  You can
// wrap *http.Client with layers of Doers to form a stack of client-side
// middleware.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Sending

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Nougat is an HTTP Request builder and sender.
type Nougat struct {
	// http Client for doing requests
	httpClient Doer
	// HTTP method (GET, POST, etc.)
	method string
	// raw url string for requests
	rawURL string
	// stores key-values pairs to add to request's Headers
	header http.Header
	// url tagged query structs
	queryStructs []interface{}
	// body provider
	bodyProvider BodyProvider
	// response decoder
	responseDecoder ResponseDecoder
}

// New returns a new Nougat with an http DefaultClient.
func New() *Nougat {
	return &Nougat{
		httpClient:      http.DefaultClient,
		method:          "GET",
		header:          make(http.Header),
		queryStructs:    make([]interface{}, 0),
		responseDecoder: jsonDecoder{},
	}
}

// New returns a copy of a Nougat for creating a new Nougat with properties
// from a parent Nougat. For example,
//
// 	parentNougat := Nougat.New().Client(client).Base("https://api.io/")
// 	fooNougat := parentNougat.New().Get("foo/")
// 	barNougat := parentNougat.New().Get("bar/")
//
// fooNougat and barNougat will both use the same client, but send requests to
// https://api.io/foo/ and https://api.io/bar/ respectively.
//
// Note that query and body values are copied so if pointer values are used,
// mutating the original value will mutate the value within the child Nougat.
func (r *Nougat) New() *Nougat {
	// copy Headers pairs into new Header map
	headerCopy := make(http.Header)
	for k, v := range r.header {
		headerCopy[k] = v
	}
	return &Nougat{
		httpClient:      r.httpClient,
		method:          r.method,
		rawURL:          r.rawURL,
		header:          headerCopy,
		queryStructs:    append([]interface{}{}, r.queryStructs...),
		bodyProvider:    r.bodyProvider,
		responseDecoder: r.responseDecoder,
	}
}

// Http Client

// Client sets the http Client used to do requests.
// If a nil client is given, the http.DefaultClient will be used.
func (r *Nougat) Client(httpClient *http.Client) *Nougat {
	if httpClient == nil {
		return r.Doer(http.DefaultClient)
	}
	return r.Doer(httpClient)
}

// Doer sets the custom Doer implementation used to do requests.
// If a nil client is given, the http.DefaultClient will be used.
func (r *Nougat) Doer(doer Doer) *Nougat {
	if doer == nil {
		r.httpClient = http.DefaultClient
	} else {
		r.httpClient = doer
	}
	return r
}
