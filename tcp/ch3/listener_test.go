package ch3_test

import (
	"net"
	"testing"
)

// Basic net listener and what network and address that you are trying to use
// also the error handling basics and making sure to close the program to not
// create memory leaks
func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	t.Logf("bound to %q", listener.Addr())
}

