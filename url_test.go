package nougat

import "testing"

func TestBaseSetter(t *testing.T) {
	cases := []string{"http://a.io/", "http://b.io", "/path", "path", ""}
	for _, base := range cases {
		Nougat := New().Base(base)
		if Nougat.rawURL != base {
			t.Errorf("expected %s, got %s", base, Nougat.rawURL)
		}
	}
}

func TestPathSetter(t *testing.T) {
	cases := []struct {
		rawURL         string
		path           string
		expectedRawURL string
	}{
		{"http://a.io/", "foo", "http://a.io/foo"},
		{"http://a.io/", "/foo", "http://a.io/foo"},
		{"http://a.io", "foo", "http://a.io/foo"},
		{"http://a.io", "/foo", "http://a.io/foo"},
		{"http://a.io/foo/", "bar", "http://a.io/foo/bar"},
		// rawURL should end in trailing slash if it is to be Path extended
		{"http://a.io/foo", "bar", "http://a.io/bar"},
		{"http://a.io/foo", "/bar", "http://a.io/bar"},
		// path extension is absolute
		{"http://a.io", "http://b.io/", "http://b.io/"},
		{"http://a.io/", "http://b.io/", "http://b.io/"},
		{"http://a.io", "http://b.io", "http://b.io"},
		{"http://a.io/", "http://b.io", "http://b.io"},
		// empty base, empty path
		{"", "http://b.io", "http://b.io"},
		{"http://a.io", "", "http://a.io"},
		{"", "", ""},
	}
	for _, c := range cases {
		Nougat := New().Base(c.rawURL).Path(c.path)
		if Nougat.rawURL != c.expectedRawURL {
			t.Errorf("expected %s, got %s", c.expectedRawURL, Nougat.rawURL)
		}
	}
}

func TestQueryStructSetter(t *testing.T) {
	cases := []struct {
		Nougat          *Nougat
		expectedStructs []interface{}
	}{
		{New(), []interface{}{}},
		{New().QueryStruct(nil), []interface{}{}},
		{New().QueryStruct(paramsA), []interface{}{paramsA}},
		{New().QueryStruct(paramsA).QueryStruct(paramsA), []interface{}{paramsA, paramsA}},
		{New().QueryStruct(paramsA).QueryStruct(paramsB), []interface{}{paramsA, paramsB}},
		{New().QueryStruct(paramsA).New(), []interface{}{paramsA}},
		{New().QueryStruct(paramsA).New().QueryStruct(paramsB), []interface{}{paramsA, paramsB}},
	}

	for _, c := range cases {
		if count := len(c.Nougat.queryStructs); count != len(c.expectedStructs) {
			t.Errorf("expected length %d, got %d", len(c.expectedStructs), count)
		}
	check:
		for _, expected := range c.expectedStructs {
			for _, param := range c.Nougat.queryStructs {
				if param == expected {
					continue check
				}
			}
			t.Errorf("expected to find %v in %v", expected, c.Nougat.queryStructs)
		}
	}
}
