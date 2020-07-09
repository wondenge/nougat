package nougat

import (
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"testing"
)

func TestReceive_success_nonDefaultDecoder(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/foo/submit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		data := ` <response>
                        <text>Some text</text>
			<favorite_count>24</favorite_count>
			<temperature>10.5</temperature>
		</response>`
		fmt.Fprintf(w, xml.Header)
		fmt.Fprintf(w, data)
	})

	endpoint := New().Client(client).Base("http://example.com/").Path("foo/").Post("submit")

	model := new(FakeModel)
	apiError := new(APIError)
	resp, err := endpoint.New().ResponseDecoder(xmlResponseDecoder{}).Receive(model, apiError)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected %d, got %d", 200, resp.StatusCode)
	}
	expectedModel := &FakeModel{Text: "Some text", FavoriteCount: 24, Temperature: 10.5}
	if !reflect.DeepEqual(expectedModel, model) {
		t.Errorf("expected %v, got %v", expectedModel, model)
	}
	expectedAPIError := &APIError{}
	if !reflect.DeepEqual(expectedAPIError, apiError) {
		t.Errorf("failureV should be zero valued, exepcted %v, got %v", expectedAPIError, apiError)
	}
}

func TestReceive_success(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/foo/submit", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertQuery(t, map[string]string{"kind_name": "vanilla", "count": "11"}, r)
		assertPostForm(t, map[string]string{"kind_name": "vanilla", "count": "11"}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"text": "Some text", "favorite_count": 24}`)
	})

	endpoint := New().Client(client).Base("http://example.com/").Path("foo/").Post("submit")
	// encode url-tagged struct in query params and as post body for testing purposes
	params := FakeParams{KindName: "vanilla", Count: 11}
	model := new(FakeModel)
	apiError := new(APIError)
	resp, err := endpoint.New().QueryStruct(params).BodyForm(params).Receive(model, apiError)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected %d, got %d", 200, resp.StatusCode)
	}
	expectedModel := &FakeModel{Text: "Some text", FavoriteCount: 24}
	if !reflect.DeepEqual(expectedModel, model) {
		t.Errorf("expected %v, got %v", expectedModel, model)
	}
	expectedAPIError := &APIError{}
	if !reflect.DeepEqual(expectedAPIError, apiError) {
		t.Errorf("failureV should be zero valued, exepcted %v, got %v", expectedAPIError, apiError)
	}
}

func TestReceive_failure(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/foo/submit", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertQuery(t, map[string]string{"kind_name": "vanilla", "count": "11"}, r)
		assertPostForm(t, map[string]string{"kind_name": "vanilla", "count": "11"}, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(429)
		fmt.Fprintf(w, `{"message": "Rate limit exceeded", "code": 88}`)
	})

	endpoint := New().Client(client).Base("http://example.com/").Path("foo/").Post("submit")
	// encode url-tagged struct in query params and as post body for testing purposes
	params := FakeParams{KindName: "vanilla", Count: 11}
	model := new(FakeModel)
	apiError := new(APIError)
	resp, err := endpoint.New().QueryStruct(params).BodyForm(params).Receive(model, apiError)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 429 {
		t.Errorf("expected %d, got %d", 429, resp.StatusCode)
	}
	expectedAPIError := &APIError{Message: "Rate limit exceeded", Code: 88}
	if !reflect.DeepEqual(expectedAPIError, apiError) {
		t.Errorf("expected %v, got %v", expectedAPIError, apiError)
	}
	expectedModel := &FakeModel{}
	if !reflect.DeepEqual(expectedModel, model) {
		t.Errorf("successV should not be zero valued, expected %v, got %v", expectedModel, model)
	}
}

func TestReceive_noContent(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/foo/submit", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "HEAD", r)
		w.WriteHeader(204)
	})

	endpoint := New().Client(client).Base("http://example.com/").Path("foo/").Head("submit")
	resp, err := endpoint.New().Receive(nil, nil)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected %d, got %d", 204, resp.StatusCode)
	}
}

func TestReceive_errorCreatingRequest(t *testing.T) {
	expectedErr := errors.New("json: unsupported value: +Inf")
	resp, err := New().BodyJSON(FakeModel{Temperature: math.Inf(1)}).Receive(nil, nil)
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
	if resp != nil {
		t.Errorf("expected nil resp, got %v", resp)
	}
}