package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestSimpleHTTPServer(t *testing.T) {
	srv := &http.Server{
		Addr:              "127.0.0.1:8081",
		Handler:           http.TimeoutHandler(handlers.DefaultHandler(), 2*time.Second, ""),
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: time.Minute, // has no effect on the time to read the request body
	}
  // we can sperately manage the time it takes to read the request body and
  // send the response using middleware or handlers. This will give us the
  // greatest control over request/reponse durations per resource
	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}

  // also for tls it is not required but common to change port 
  // resemblance of 443 - 3443 - 8443
	// go func() {
	// 	err := srv.ServeTLS(l, "cert.pem", "key.pem")
	// 	if err != http.ErrServerClosed {
	// 		t.Error(err)
	// 	}
	// }()

	go func() {
		err := srv.Serve(l)
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	testCases := []struct {
		method   string
		body     io.Reader
		code     int
		response string
	}{
		{http.MethodGet, nil, http.StatusOK, "Hello, friend!"},
		{http.MethodPost, bytes.NewBufferString("<world>"), http.StatusOK, "Hello, &lt;world&gt;!"},
		{http.MethodHead, nil, http.StatusMethodNotAllowed, ""},
	}
	client := new(http.Client)
	path := fmt.Sprintf("http://%s/", srv.Addr)

	for i, c := range testCases {
		r, err := http.NewRequest(c.method, path, c.body)
		if err != nil {
			t.Errorf("%d: %v\n", i, err)
			continue
		}
		resp, err := client.Do(r)
		if err != nil {
			t.Errorf("%d: %v\n", i, err)
			continue
		}
		if resp.StatusCode != c.code {
			t.Errorf("%d: unexpected status code: %q\n", i, resp.Status)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("%d: %v\n", i, err)
			continue
		}
		_ = resp.Body.Close()
		if c.response != string(b) {
			t.Errorf("%d: expected %q; actual %q\n", i, c.response, b)
		}
	}
	if err := srv.Close(); err != nil {
		t.Fatal(err)
	}
}
