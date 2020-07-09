package nougat

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestBodySetter(t *testing.T) {
	fakeInput := ioutil.NopCloser(strings.NewReader("test"))
	fakeBodyProvider := bodyProvider{body: fakeInput}

	cases := []struct {
		initial  BodyProvider
		input    io.Reader
		expected BodyProvider
	}{
		// nil body is overriden by a set body
		{nil, fakeInput, fakeBodyProvider},
		// initial body is not overriden by nil body
		{fakeBodyProvider, nil, fakeBodyProvider},
		// nil body is returned unaltered
		{nil, nil, nil},
	}
	for _, c := range cases {
		Nougat := New()
		Nougat.bodyProvider = c.initial
		Nougat.Body(c.input)
		if Nougat.bodyProvider != c.expected {
			t.Errorf("expected %v, got %v", c.expected, Nougat.bodyProvider)
		}
	}
}

func TestBodyJSONSetter(t *testing.T) {
	fakeModel := &FakeModel{}
	fakeBodyProvider := jsonBodyProvider{payload: fakeModel}

	cases := []struct {
		initial  BodyProvider
		input    interface{}
		expected BodyProvider
	}{
		// json tagged struct is set as bodyJSON
		{nil, fakeModel, fakeBodyProvider},
		// nil argument to bodyJSON does not replace existing bodyJSON
		{fakeBodyProvider, nil, fakeBodyProvider},
		// nil bodyJSON remains nil
		{nil, nil, nil},
	}
	for _, c := range cases {
		Nougat := New()
		Nougat.bodyProvider = c.initial
		Nougat.BodyJSON(c.input)
		if Nougat.bodyProvider != c.expected {
			t.Errorf("expected %v, got %v", c.expected, Nougat.bodyProvider)
		}
		// Header Content-Type should be application/json if bodyJSON arg was non-nil
		if c.input != nil && Nougat.header.Get(contentType) != jsonContentType {
			t.Errorf("Incorrect or missing header, expected %s, got %s", jsonContentType, Nougat.header.Get(contentType))
		} else if c.input == nil && Nougat.header.Get(contentType) != "" {
			t.Errorf("did not expect a Content-Type header, got %s", Nougat.header.Get(contentType))
		}
	}
}

func TestBodyFormSetter(t *testing.T) {
	fakeParams := FakeParams{KindName: "recent", Count: 25}
	fakeBodyProvider := formBodyProvider{payload: fakeParams}

	cases := []struct {
		initial  BodyProvider
		input    interface{}
		expected BodyProvider
	}{
		// url tagged struct is set as bodyStruct
		{nil, paramsB, fakeBodyProvider},
		// nil argument to bodyStruct does not replace existing bodyStruct
		{fakeBodyProvider, nil, fakeBodyProvider},
		// nil bodyStruct remains nil
		{nil, nil, nil},
	}
	for _, c := range cases {
		Nougat := New()
		Nougat.bodyProvider = c.initial
		Nougat.BodyForm(c.input)
		if Nougat.bodyProvider != c.expected {
			t.Errorf("expected %v, got %v", c.expected, Nougat.bodyProvider)
		}
		// Content-Type should be application/x-www-form-urlencoded if bodyStruct was non-nil
		if c.input != nil && Nougat.header.Get(contentType) != formContentType {
			t.Errorf("Incorrect or missing header, expected %s, got %s", formContentType, Nougat.header.Get(contentType))
		} else if c.input == nil && Nougat.header.Get(contentType) != "" {
			t.Errorf("did not expect a Content-Type header, got %s", Nougat.header.Get(contentType))
		}
	}
}
