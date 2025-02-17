package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHeadTime(t *testing.T) {
  resp, err := http.Head("https://time.gov")
  if err != nil {
    t.Fatal(err)
  }
  _ = resp.Body.Close()
  now := time.Now().Round(time.Second)
  date := resp.Header.Get("Date")
  if date == "" {
    t.Fatal("no date header received from time.gov")
  }
  dt, err := time.Parse(time.RFC1123, date)
  if err != nil {
    t.Fatal(err)
  }
  t.Logf("time.gov: %s (skew %s)\n", dt, now.Sub(dt))
}

func blockIndefinitely(w http.ResponseWriter, r *http.Request) {
  select {}
}

func TestBlockIndefinitely(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))
  _, _ = http.Get(ts.URL)
  t.Fatal("client did not indefinitely block")
}
