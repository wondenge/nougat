package nougat

import (
	"net/url"
	"reflect"
	"testing"
)

func TestAddHeader(t *testing.T) {
	cases := []struct {
		Nougat         *Nougat
		expectedHeader map[string][]string
	}{
		{New().Add("authorization", "OAuth key=\"value\""), map[string][]string{"Authorization": []string{"OAuth key=\"value\""}}},
		// header keys should be canonicalized
		{New().Add("content-tYPE", "application/json").Add("User-AGENT", "Nougat"), map[string][]string{"Content-Type": []string{"application/json"}, "User-Agent": []string{"Nougat"}}},
		// values for existing keys should be appended
		{New().Add("A", "B").Add("a", "c"), map[string][]string{"A": []string{"B", "c"}}},
		// Add should add to values for keys added by parent Nougats
		{New().Add("A", "B").Add("a", "c").New(), map[string][]string{"A": []string{"B", "c"}}},
		{New().Add("A", "B").New().Add("a", "c"), map[string][]string{"A": []string{"B", "c"}}},
	}
	for _, c := range cases {
		// type conversion from header to alias'd map for deep equality comparison
		headerMap := map[string][]string(c.Nougat.header)
		if !reflect.DeepEqual(c.expectedHeader, headerMap) {
			t.Errorf("not DeepEqual: expected %v, got %v", c.expectedHeader, headerMap)
		}
	}
}

func TestSetHeader(t *testing.T) {
	cases := []struct {
		Nougat         *Nougat
		expectedHeader map[string][]string
	}{
		// should replace existing values associated with key
		{New().Add("A", "B").Set("a", "c"), map[string][]string{"A": []string{"c"}}},
		{New().Set("content-type", "A").Set("Content-Type", "B"), map[string][]string{"Content-Type": []string{"B"}}},
		// Set should replace values received by copying parent Nougats
		{New().Set("A", "B").Add("a", "c").New(), map[string][]string{"A": []string{"B", "c"}}},
		{New().Add("A", "B").New().Set("a", "c"), map[string][]string{"A": []string{"c"}}},
	}
	for _, c := range cases {
		// type conversion from Header to alias'd map for deep equality comparison
		headerMap := map[string][]string(c.Nougat.header)
		if !reflect.DeepEqual(c.expectedHeader, headerMap) {
			t.Errorf("not DeepEqual: expected %v, got %v", c.expectedHeader, headerMap)
		}
	}
}

func TestBasicAuth(t *testing.T) {
	cases := []struct {
		Nougat       *Nougat
		expectedAuth []string
	}{
		// basic auth: username & password
		{New().SetBasicAuth("Aladdin", "open sesame"), []string{"Aladdin", "open sesame"}},
		// empty username
		{New().SetBasicAuth("", "secret"), []string{"", "secret"}},
		// empty password
		{New().SetBasicAuth("admin", ""), []string{"admin", ""}},
	}
	for _, c := range cases {
		req, err := c.Nougat.Request()
		if err != nil {
			t.Errorf("unexpected error when building Request with .SetBasicAuth()")
		}
		username, password, ok := req.BasicAuth()
		if !ok {
			t.Errorf("basic auth missing when expected")
		}
		auth := []string{username, password}
		if !reflect.DeepEqual(c.expectedAuth, auth) {
			t.Errorf("not DeepEqual: expected %v, got %v", c.expectedAuth, auth)
		}
	}
}


func TestAddQueryStructs(t *testing.T) {
	cases := []struct {
		rawurl       string
		queryStructs []interface{}
		expected     string
	}{
		{"http://a.io", []interface{}{}, "http://a.io"},
		{"http://a.io", []interface{}{paramsA}, "http://a.io?limit=30"},
		{"http://a.io", []interface{}{paramsA, paramsA}, "http://a.io?limit=30&limit=30"},
		{"http://a.io", []interface{}{paramsA, paramsB}, "http://a.io?count=25&kind_name=recent&limit=30"},
		// don't blow away query values on the rawURL (parsed into RawQuery)
		{"http://a.io?initial=7", []interface{}{paramsA}, "http://a.io?initial=7&limit=30"},
	}
	for _, c := range cases {
		reqURL, _ := url.Parse(c.rawurl)
		addQueryStructs(reqURL, c.queryStructs)
		if reqURL.String() != c.expected {
			t.Errorf("expected %s, got %s", c.expected, reqURL.String())
		}
	}
}


func TestRequest_headers(t *testing.T) {
	cases := []struct {
		Nougat         *Nougat
		expectedHeader map[string][]string
	}{
		{New().Add("authorization", "OAuth key=\"value\""), map[string][]string{"Authorization": []string{"OAuth key=\"value\""}}},
		// header keys should be canonicalized
		{New().Add("content-tYPE", "application/json").Add("User-AGENT", "Nougat"), map[string][]string{"Content-Type": []string{"application/json"}, "User-Agent": []string{"Nougat"}}},
		// values for existing keys should be appended
		{New().Add("A", "B").Add("a", "c"), map[string][]string{"A": []string{"B", "c"}}},
		// Add should add to values for keys added by parent Nougats
		{New().Add("A", "B").Add("a", "c").New(), map[string][]string{"A": []string{"B", "c"}}},
		{New().Add("A", "B").New().Add("a", "c"), map[string][]string{"A": []string{"B", "c"}}},
		// Add and Set
		{New().Add("A", "B").Set("a", "c"), map[string][]string{"A": []string{"c"}}},
		{New().Set("content-type", "A").Set("Content-Type", "B"), map[string][]string{"Content-Type": []string{"B"}}},
		// Set should replace values received by copying parent Nougats
		{New().Set("A", "B").Add("a", "c").New(), map[string][]string{"A": []string{"B", "c"}}},
		{New().Add("A", "B").New().Set("a", "c"), map[string][]string{"A": []string{"c"}}},
	}
	for _, c := range cases {
		req, _ := c.Nougat.Request()
		// type conversion from Header to alias'd map for deep equality comparison
		headerMap := map[string][]string(req.Header)
		if !reflect.DeepEqual(c.expectedHeader, headerMap) {
			t.Errorf("not DeepEqual: expected %v, got %v", c.expectedHeader, headerMap)
		}
	}
}
