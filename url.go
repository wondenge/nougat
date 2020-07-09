package nougat

import "net/url"

// Base sets the rawURL.
// If you intend to extend the url with Path, baseUrl should be specified
// with a trailing slash.
func (r *Nougat) Base(rawURL string) *Nougat {
	r.rawURL = rawURL
	return r
}

// Path extends the rawURL with the given path by resolving the reference
// to an absolute URL.
// If parsing errors occur, the rawURL is left unmodified.
func (r *Nougat) Path(path string) *Nougat {
	baseURL, baseErr := url.Parse(r.rawURL)
	pathURL, pathErr := url.Parse(path)
	if baseErr == nil && pathErr == nil {
		r.rawURL = baseURL.ResolveReference(pathURL).String()
		return r
	}
	return r
}

// QueryStruct appends the queryStruct to the Nougat's queryStructs.
// The value pointed to by each queryStruct will be encoded as url query
// parameters on new requests (see Request()).
// The queryStruct argument should be a pointer to a url tagged struct.
// See https://godoc.org/github.com/google/go-querystring/query for details.
func (r *Nougat) QueryStruct(queryStruct interface{}) *Nougat {
	if queryStruct != nil {
		r.queryStructs = append(r.queryStructs, queryStruct)
	}
	return r
}
