package nougat

import (
	"bytes"
	"errors"
	"math"
	"strings"
	"testing"
)

func TestRequest_urlAndMethod(t *testing.T) {
	cases := []struct {
		Nougat         *Nougat
		expectedMethod string
		expectedURL    string
		expectedErr    error
	}{
		{New().Base("http://a.io"), "GET", "http://a.io", nil},
		{New().Path("http://a.io"), "GET", "http://a.io", nil},
		{New().Get("http://a.io"), "GET", "http://a.io", nil},
		{New().Put("http://a.io"), "PUT", "http://a.io", nil},
		{New().Base("http://a.io/").Path("foo"), "GET", "http://a.io/foo", nil},
		{New().Base("http://a.io/").Post("foo"), "POST", "http://a.io/foo", nil},
		// if relative path is an absolute url, base is ignored
		{New().Base("http://a.io").Path("http://b.io"), "GET", "http://b.io", nil},
		{New().Path("http://a.io").Path("http://b.io"), "GET", "http://b.io", nil},
		// last method setter takes priority
		{New().Get("http://b.io").Post("http://a.io"), "POST", "http://a.io", nil},
		{New().Post("http://a.io/").Put("foo/").Delete("bar"), "DELETE", "http://a.io/foo/bar", nil},
		// last Base setter takes priority
		{New().Base("http://a.io").Base("http://b.io"), "GET", "http://b.io", nil},
		// Path setters are additive
		{New().Base("http://a.io/").Path("foo/").Path("bar"), "GET", "http://a.io/foo/bar", nil},
		{New().Path("http://a.io/").Path("foo/").Path("bar"), "GET", "http://a.io/foo/bar", nil},
		// removes extra '/' between base and ref url
		{New().Base("http://a.io/").Get("/foo"), "GET", "http://a.io/foo", nil},
	}
	for _, c := range cases {
		req, err := c.Nougat.Request()
		if err != c.expectedErr {
			t.Errorf("expected error %v, got %v for %+v", c.expectedErr, err, c.Nougat)
		}
		if req.URL.String() != c.expectedURL {
			t.Errorf("expected url %s, got %s for %+v", c.expectedURL, req.URL.String(), c.Nougat)
		}
		if req.Method != c.expectedMethod {
			t.Errorf("expected method %s, got %s for %+v", c.expectedMethod, req.Method, c.Nougat)
		}
	}
}

func TestRequest_queryStructs(t *testing.T) {
	cases := []struct {
		Nougat      *Nougat
		expectedURL string
	}{
		{New().Base("http://a.io").QueryStruct(paramsA), "http://a.io?limit=30"},
		{New().Base("http://a.io").QueryStruct(paramsA).QueryStruct(paramsB), "http://a.io?count=25&kind_name=recent&limit=30"},
		{New().Base("http://a.io/").Path("foo?path=yes").QueryStruct(paramsA), "http://a.io/foo?limit=30&path=yes"},
		{New().Base("http://a.io").QueryStruct(paramsA).New(), "http://a.io?limit=30"},
		{New().Base("http://a.io").QueryStruct(paramsA).New().QueryStruct(paramsB), "http://a.io?count=25&kind_name=recent&limit=30"},
	}
	for _, c := range cases {
		req, _ := c.Nougat.Request()
		if req.URL.String() != c.expectedURL {
			t.Errorf("expected url %s, got %s for %+v", c.expectedURL, req.URL.String(), c.Nougat)
		}
	}
}

func TestRequest_body(t *testing.T) {
	cases := []struct {
		Nougat              *Nougat
		expectedBody        string // expected Body io.Reader as a string
		expectedContentType string
	}{
		// BodyJSON
		{New().BodyJSON(modelA), "{\"text\":\"note\",\"favorite_count\":12}\n", jsonContentType},
		{New().BodyJSON(&modelA), "{\"text\":\"note\",\"favorite_count\":12}\n", jsonContentType},
		{New().BodyJSON(&FakeModel{}), "{}\n", jsonContentType},
		{New().BodyJSON(FakeModel{}), "{}\n", jsonContentType},
		// BodyJSON overrides existing values
		{New().BodyJSON(&FakeModel{}).BodyJSON(&FakeModel{Text: "msg"}), "{\"text\":\"msg\"}\n", jsonContentType},
		// BodyForm
		{New().BodyForm(paramsA), "limit=30", formContentType},
		{New().BodyForm(paramsB), "count=25&kind_name=recent", formContentType},
		{New().BodyForm(&paramsB), "count=25&kind_name=recent", formContentType},
		// BodyForm overrides existing values
		{New().BodyForm(paramsA).New().BodyForm(paramsB), "count=25&kind_name=recent", formContentType},
		// Mixture of BodyJSON and BodyForm prefers body setter called last with a non-nil argument
		{New().BodyForm(paramsB).New().BodyJSON(modelA), "{\"text\":\"note\",\"favorite_count\":12}\n", jsonContentType},
		{New().BodyJSON(modelA).New().BodyForm(paramsB), "count=25&kind_name=recent", formContentType},
		{New().BodyForm(paramsB).New().BodyJSON(nil), "count=25&kind_name=recent", formContentType},
		{New().BodyJSON(modelA).New().BodyForm(nil), "{\"text\":\"note\",\"favorite_count\":12}\n", jsonContentType},
		// Body
		{New().Body(strings.NewReader("this-is-a-test")), "this-is-a-test", ""},
		{New().Body(strings.NewReader("a")).Body(strings.NewReader("b")), "b", ""},
	}
	for _, c := range cases {
		req, _ := c.Nougat.Request()
		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		// req.Body should have contained the expectedBody string
		if value := buf.String(); value != c.expectedBody {
			t.Errorf("expected Request.Body %s, got %s", c.expectedBody, value)
		}
		// Header Content-Type should be expectedContentType ("" means no contentType expected)
		if actualHeader := req.Header.Get(contentType); actualHeader != c.expectedContentType && c.expectedContentType != "" {
			t.Errorf("Incorrect or missing header, expected %s, got %s", c.expectedContentType, actualHeader)
		}
	}
}

func TestRequest_bodyNoData(t *testing.T) {
	// test that Body is left nil when no bodyJSON or bodyStruct set
	Nougats := []*Nougat{
		New(),
		New().BodyJSON(nil),
		New().BodyForm(nil),
	}
	for _, Nougat := range Nougats {
		req, _ := Nougat.Request()
		if req.Body != nil {
			t.Errorf("expected nil Request.Body, got %v", req.Body)
		}
		// Header Content-Type should not be set when bodyJSON argument was nil or never called
		if actualHeader := req.Header.Get(contentType); actualHeader != "" {
			t.Errorf("did not expect a Content-Type header, got %s", actualHeader)
		}
	}
}

func TestRequest_bodyEncodeErrors(t *testing.T) {
	cases := []struct {
		Nougat      *Nougat
		expectedErr error
	}{
		// check that Encode errors are propagated, illegal JSON field
		{New().BodyJSON(FakeModel{Temperature: math.Inf(1)}), errors.New("json: unsupported value: +Inf")},
	}
	for _, c := range cases {
		req, err := c.Nougat.Request()
		if err == nil || err.Error() != c.expectedErr.Error() {
			t.Errorf("expected error %v, got %v", c.expectedErr, err)
		}
		if req != nil {
			t.Errorf("expected nil Request, got %+v", req)
		}
	}
}
