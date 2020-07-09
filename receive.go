package nougat

import "net/http"

// ReceiveSuccess creates a new HTTP request and returns the response. Success
// responses (2XX) are JSON decoded into the value pointed to by successV.
// Any error creating the request, sending it, or decoding a 2XX response
// is returned.
func (r *Nougat) ReceiveSuccess(successV interface{}) (*http.Response, error) {
	return r.Receive(successV, nil)
}

// Receive creates a new HTTP request and returns the response. Success
// responses (2XX) are JSON decoded into the value pointed to by successV and
// other responses are JSON decoded into the value pointed to by failureV.
// If the status code of response is 204(no content), decoding is skipped.
// Any error creating the request, sending it, or decoding the response is
// returned.
// Receive is shorthand for calling Request and Do.
func (r *Nougat) Receive(successV, failureV interface{}) (*http.Response, error) {
	req, err := r.Request()
	if err != nil {
		return nil, err
	}
	return r.Do(req, successV, failureV)
}
