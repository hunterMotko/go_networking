package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func blockIndefinitley(w http.ResponseWriter, r *http.Request) {
  select {}
}

func TestBlockIndefinitley(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))
  _, _ = http.Get(ts.URL)
  t.Fatal("client did not indefinitley block")
}

func TestBlockIndefinitleyWithTimeout(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  // when using context this might timeout before the body is read
  // timer := time.AfterFunc(5*time.second, cancel)
  // Make the HTTP request, Read the response headers
  req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
  if err != nil {
    t.Fatal(err)
  }
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    if !errors.Is(err, context.DeadlineExceeded) {
      t.Fatal(err)
    }
    return
  }
  // Add 5 more seconds before reading the response body
  // timer.Reset(5*time.Second)
  _ = resp.Body.Close()
}
