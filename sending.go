package nougat

import (
	"net/http"
)

const (
	contentType = "Content-Type"
)

// Sending

// ResponseDecoder sets the Nougat's response decoder.
func (r *Nougat) ResponseDecoder(decoder ResponseDecoder) *Nougat {
	if decoder == nil {
		return r
	}
	r.responseDecoder = decoder
	return r
}

// decodeResponse decodes response Body into the value pointed to by successV
// if the response is a success (2XX) or into the value pointed to by failureV
// otherwise. If the successV or failureV argument to decode into is nil,
// decoding is skipped.
// Caller is responsible for closing the resp.Body.
func decodeResponse(resp *http.Response, decoder ResponseDecoder, successV, failureV interface{}) error {
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if successV != nil {
			return decoder.Decode(resp, successV)
		}
	} else {
		if failureV != nil {
			return decoder.Decode(resp, failureV)
		}
	}
	return nil
}
