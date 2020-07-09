package nougat

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDo_onSuccess(t *testing.T) {
	const expectedText = "Some text"
	const expectedFavoriteCount int64 = 24

	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"text": "Some text", "favorite_count": 24}`)
	})

	Nougat := New().Client(client)
	req, _ := http.NewRequest("GET", "http://example.com/success", nil)

	model := new(FakeModel)
	apiError := new(APIError)
	resp, err := Nougat.Do(req, model, apiError)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected %d, got %d", 200, resp.StatusCode)
	}
	if model.Text != expectedText {
		t.Errorf("expected %s, got %s", expectedText, model.Text)
	}
	if model.FavoriteCount != expectedFavoriteCount {
		t.Errorf("expected %d, got %d", expectedFavoriteCount, model.FavoriteCount)
	}
}

func TestDo_onSuccessWithNilValue(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"text": "Some text", "favorite_count": 24}`)
	})

	Nougat := New().Client(client)
	req, _ := http.NewRequest("GET", "http://example.com/success", nil)

	apiError := new(APIError)
	resp, err := Nougat.Do(req, nil, apiError)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected %d, got %d", 200, resp.StatusCode)
	}
	expected := &APIError{}
	if !reflect.DeepEqual(expected, apiError) {
		t.Errorf("failureV should not be populated, exepcted %v, got %v", expected, apiError)
	}
}

func TestDo_noContent(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/nocontent", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	Nougat := New().Client(client)
	req, _ := http.NewRequest("DELETE", "http://example.com/nocontent", nil)

	model := new(FakeModel)
	apiError := new(APIError)
	resp, err := Nougat.Do(req, model, apiError)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected %d, got %d", 204, resp.StatusCode)
	}
	expectedModel := &FakeModel{}
	if !reflect.DeepEqual(expectedModel, model) {
		t.Errorf("successV should not be populated, exepcted %v, got %v", expectedModel, model)
	}
	expectedAPIError := &APIError{}
	if !reflect.DeepEqual(expectedAPIError, apiError) {
		t.Errorf("failureV should not be populated, exepcted %v, got %v", expectedAPIError, apiError)
	}
}

func TestDo_onFailure(t *testing.T) {
	const expectedMessage = "Invalid argument"
	const expectedCode int = 215

	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/failure", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Invalid argument", "code": 215}`)
	})

	Nougat := New().Client(client)
	req, _ := http.NewRequest("GET", "http://example.com/failure", nil)

	model := new(FakeModel)
	apiError := new(APIError)
	resp, err := Nougat.Do(req, model, apiError)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected %d, got %d", 400, resp.StatusCode)
	}
	if apiError.Message != expectedMessage {
		t.Errorf("expected %s, got %s", expectedMessage, apiError.Message)
	}
	if apiError.Code != expectedCode {
		t.Errorf("expected %d, got %d", expectedCode, apiError.Code)
	}
}

func TestDo_onFailureWithNilValue(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/failure", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(420)
		fmt.Fprintf(w, `{"message": "Enhance your calm", "code": 88}`)
	})

	Nougat := New().Client(client)
	req, _ := http.NewRequest("GET", "http://example.com/failure", nil)

	model := new(FakeModel)
	resp, err := Nougat.Do(req, model, nil)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if resp.StatusCode != 420 {
		t.Errorf("expected %d, got %d", 420, resp.StatusCode)
	}
	expected := &FakeModel{}
	if !reflect.DeepEqual(expected, model) {
		t.Errorf("successV should not be populated, exepcted %v, got %v", expected, model)
	}
}
