package nougat

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)


type FakeParams struct {
	KindName string `url:"kind_name"`
	Count    int    `url:"count"`
}

// Url-tagged query struct
var paramsA = struct {
	Limit int `url:"limit"`
}{
	30,
}
var paramsB = FakeParams{KindName: "recent", Count: 25}

// Json/XML-tagged model struct
type FakeModel struct {
	Text          string  `json:"text,omitempty" xml:"text"`
	FavoriteCount int64   `json:"favorite_count,omitempty" xml:"favorite_count"`
	Temperature   float64 `json:"temperature,omitempty" xml:"temperature"`
}

var modelA = FakeModel{Text: "note", FavoriteCount: 12}

// Non-Json response decoder
type xmlResponseDecoder struct{}

func (d xmlResponseDecoder) Decode(resp *http.Response, v interface{}) error {
	return xml.NewDecoder(resp.Body).Decode(v)
}

func TestNew(t *testing.T) {
	Nougat := New()
	if Nougat.httpClient != http.DefaultClient {
		t.Errorf("expected %v, got %v", http.DefaultClient, Nougat.httpClient)
	}
	if Nougat.header == nil {
		t.Errorf("Header map not initialized with make")
	}
	if Nougat.queryStructs == nil {
		t.Errorf("queryStructs not initialized with make")
	}
}

func TestNougatNew(t *testing.T) {
	fakeBodyProvider := jsonBodyProvider{FakeModel{}}

	cases := []*Nougat{
		&Nougat{httpClient: &http.Client{}, method: "GET", rawURL: "http://example.com"},
		&Nougat{httpClient: nil, method: "", rawURL: "http://example.com"},
		&Nougat{queryStructs: make([]interface{}, 0)},
		&Nougat{queryStructs: []interface{}{paramsA}},
		&Nougat{queryStructs: []interface{}{paramsA, paramsB}},
		&Nougat{bodyProvider: fakeBodyProvider},
		&Nougat{bodyProvider: fakeBodyProvider},
		&Nougat{bodyProvider: nil},
		New().Add("Content-Type", "application/json"),
		New().Add("A", "B").Add("a", "c").New(),
		New().Add("A", "B").New().Add("a", "c"),
		New().BodyForm(paramsB),
		New().BodyForm(paramsB).New(),
	}
	for _, Nougat := range cases {
		child := Nougat.New()
		if child.httpClient != Nougat.httpClient {
			t.Errorf("expected %v, got %v", Nougat.httpClient, child.httpClient)
		}
		if child.method != Nougat.method {
			t.Errorf("expected %s, got %s", Nougat.method, child.method)
		}
		if child.rawURL != Nougat.rawURL {
			t.Errorf("expected %s, got %s", Nougat.rawURL, child.rawURL)
		}
		// Header should be a copy of parent Nougat header. For example, calling
		// baseNougat.Add("k","v") should not mutate previously created child Nougats
		if Nougat.header != nil {
			// struct literal cases don't init Header in usual way, skip header check
			if !reflect.DeepEqual(Nougat.header, child.header) {
				t.Errorf("not DeepEqual: expected %v, got %v", Nougat.header, child.header)
			}
			Nougat.header.Add("K", "V")
			if child.header.Get("K") != "" {
				t.Errorf("child.header was a reference to original map, should be copy")
			}
		}
		// queryStruct slice should be a new slice with a copy of the contents
		if len(Nougat.queryStructs) > 0 {
			// mutating one slice should not mutate the other
			child.queryStructs[0] = nil
			if Nougat.queryStructs[0] == nil {
				t.Errorf("child.queryStructs was a re-slice, expected slice with copied contents")
			}
		}
		// body should be copied
		if child.bodyProvider != Nougat.bodyProvider {
			t.Errorf("expected %v, got %v", Nougat.bodyProvider, child.bodyProvider)
		}
	}
}

func TestClientSetter(t *testing.T) {
	developerClient := &http.Client{}
	cases := []struct {
		input    *http.Client
		expected *http.Client
	}{
		{nil, http.DefaultClient},
		{developerClient, developerClient},
	}
	for _, c := range cases {
		Nougat := New()
		Nougat.Client(c.input)
		if Nougat.httpClient != c.expected {
			t.Errorf("input %v, expected %v, got %v", c.input, c.expected, Nougat.httpClient)
		}
	}
}

func TestDoerSetter(t *testing.T) {
	developerClient := &http.Client{}
	cases := []struct {
		input    Doer
		expected Doer
	}{
		{nil, http.DefaultClient},
		{developerClient, developerClient},
	}
	for _, c := range cases {
		Nougat := New()
		Nougat.Doer(c.input)
		if Nougat.httpClient != c.expected {
			t.Errorf("input %v, expected %v, got %v", c.input, c.expected, Nougat.httpClient)
		}
	}
}

// Testing Utils

// testServer returns an http Client, ServeMux, and Server. The client proxies
// requests to the server and handlers can be registered on the mux to handle
// requests. The caller must close the test server.
func testServer() (*http.Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	client := &http.Client{Transport: transport}
	return client, mux, server
}

func assertMethod(t *testing.T, expectedMethod string, req *http.Request) {
	if actualMethod := req.Method; actualMethod != expectedMethod {
		t.Errorf("expected method %s, got %s", expectedMethod, actualMethod)
	}
}

// assertQuery tests that the Request has the expected url query key/val pairs
func assertQuery(t *testing.T, expected map[string]string, req *http.Request) {
	queryValues := req.URL.Query() // net/url Values is a map[string][]string
	expectedValues := url.Values{}
	for key, value := range expected {
		expectedValues.Add(key, value)
	}
	if !reflect.DeepEqual(expectedValues, queryValues) {
		t.Errorf("expected parameters %v, got %v", expected, req.URL.RawQuery)
	}
}

// assertPostForm tests that the Request has the expected key values pairs url
// encoded in its Body
func assertPostForm(t *testing.T, expected map[string]string, req *http.Request) {
	req.ParseForm() // parses request Body to put url.Values in r.Form/r.PostForm
	expectedValues := url.Values{}
	for key, value := range expected {
		expectedValues.Add(key, value)
	}
	if !reflect.DeepEqual(expectedValues, req.PostForm) {
		t.Errorf("expected parameters %v, got %v", expected, req.PostForm)
	}
}
