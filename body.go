package nougat

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	jsonContentType = "application/json"
	formContentType = "application/x-www-form-urlencoded"
)

type (
	BodyProvider interface {

		// ContentType returns the Content-Type of the body
		ContentType() string

		// Body returns the io.Reader body.
		Body() (io.Reader, error)
	}

	// bodyProvider provides the wrapped body value as a Body for requests.
	bodyProvider struct {
		body io.Reader
	}

	// jsonBodyProvider encodes a JSON tagged struct value as a Body for requests.
	// See https://golang.org/pkg/encoding/json/#MarshalIndent for details.
	jsonBodyProvider struct {
		payload interface{}
	}

	// formBodyProvider encodes a url tagged struct value as Body for requests.
	// See https://godoc.org/github.com/google/go-querystring/query for details.
	formBodyProvider struct {
		payload interface{}
	}
)

// Body
/********************************** BODY *********************************************/

func (p bodyProvider) ContentType() string {
	return ""
}

func (p bodyProvider) Body() (io.Reader, error) {
	return p.body, nil
}

// Body sets the Nougat's body.
// The body value will be set as the Body on new requests (see Request()).
// If the provided body is also an io.Closer, the request Body will be closed by http.Client methods.
func (r *Nougat) Body(body io.Reader) *Nougat {
	if body == nil {
		return r
	}
	return r.BodyProvider(bodyProvider{body: body})
}

// BodyProvider sets the Nougat's body provider.
func (r *Nougat) BodyProvider(body BodyProvider) *Nougat {
	if body == nil {
		return r
	}
	r.bodyProvider = body

	ct := body.ContentType()
	if ct != "" {
		r.Set(contentType, ct)
	}

	return r
}

/********************************** JSON BODY *********************************************/

func (p jsonBodyProvider) ContentType() string {
	return jsonContentType
}

func (p jsonBodyProvider) Body() (io.Reader, error) {
	buf := &bytes.Buffer{}

	err := json.NewEncoder(buf).Encode(p.payload)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// BodyJSON sets the Nougat's bodyJSON.
// The value pointed to by the bodyJSON will be JSON encoded as the Body on new requests (see Request()).
// The bodyJSON argument should be a pointer to a JSON tagged struct.
// See https://golang.org/pkg/encoding/json/#MarshalIndent for details.
func (r *Nougat) BodyJSON(bodyJSON interface{}) *Nougat {
	if bodyJSON == nil {
		return r
	}
	return r.BodyProvider(jsonBodyProvider{payload: bodyJSON})
}

/********************************** FORM BODY *********************************************/
func (p formBodyProvider) ContentType() string {
	return formContentType
}

func (p formBodyProvider) Body() (io.Reader, error) {
	values, err := query.Values(p.payload)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(values.Encode()), nil
}

// BodyForm sets the Nougat's bodyForm.
// The value pointed to by the bodyForm will be url encoded as the Body on new requests (see Request()).
// The bodyForm argument should be a pointer to a url tagged struct.
// See https://godoc.org/github.com/google/go-querystring/query for details.
func (r *Nougat) BodyForm(bodyForm interface{}) *Nougat {
	if bodyForm == nil {
		return r
	}
	return r.BodyProvider(formBodyProvider{payload: bodyForm})
}
