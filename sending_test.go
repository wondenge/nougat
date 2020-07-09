package nougat

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"testing"
)

func TestReuseTcpConnections(t *testing.T) {
	var connCount int32

	ln, _ := net.Listen("tcp", ":0")
	rawURL := fmt.Sprintf("http://%s/", ln.Addr())

	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertMethod(t, "GET", r)
			fmt.Fprintf(w, `{"text": "Some text"}`)
		}),
		ConnState: func(conn net.Conn, state http.ConnState) {
			if state == http.StateNew {
				atomic.AddInt32(&connCount, 1)
			}
		},
	}

	go server.Serve(ln)

	endpoint := New().Client(http.DefaultClient).Base(rawURL).Path("foo/").Get("get")

	for i := 0; i < 10; i++ {
		resp, err := endpoint.New().Receive(nil, nil)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("expected %d, got %d", 200, resp.StatusCode)
		}
	}

	server.Shutdown(context.Background())

	if count := atomic.LoadInt32(&connCount); count != 1 {
		t.Errorf("expected 1, got %v", count)
	}
}