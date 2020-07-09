package nougat

import "testing"

func TestMethodSetters(t *testing.T) {
	cases := []struct {
		Nougat         *Nougat
		expectedMethod string
	}{
		{New().Path("http://a.io"), "GET"},
		{New().Head("http://a.io"), "HEAD"},
		{New().Get("http://a.io"), "GET"},
		{New().Post("http://a.io"), "POST"},
		{New().Put("http://a.io"), "PUT"},
		{New().Patch("http://a.io"), "PATCH"},
		{New().Delete("http://a.io"), "DELETE"},
		{New().Options("http://a.io"), "OPTIONS"},
		{New().Trace("http://a.io"), "TRACE"},
		{New().Connect("http://a.io"), "CONNECT"},
	}
	for _, c := range cases {
		if c.Nougat.method != c.expectedMethod {
			t.Errorf("expected method %s, got %s", c.expectedMethod, c.Nougat.method)
		}
	}
}
