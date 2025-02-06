package ch3_test

import (
	"net"
	"syscall"
	"testing"
	"time"
)

// probably not exactly how you are going to want to handle this problem because if you are going 
// this deep you are potentially looking for the underlying problem of why the connection
// is actually not connecting or why it is timing out.
// This is just to show a brief example of how to can use the net struct to access the contorl
// for your use case of the net dial for your program. This is just to show the control 
// over the program and the strcut useage in GO
func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
  d := net.Dialer{
    Control: func(_, address string, _ syscall.RawConn) error {
      return &net.DNSError{
        Err: "connection timed out",
        Name: address,
        Server: "127.0.0.1",
        IsTimeout: true,
        IsTemporary: true,
      }
    },
    Timeout: timeout,
  }
  return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {
  c, err := DialTimeout("tcp", "10.0.0.1:http", 5 * time.Second)
  if err == nil {
    c.Close()
    t.Fatal("connection did not timeout")
  }
  nErr, ok := err.(net.Error)
  if !ok {
    t.Fatal(err)
  }
  if !nErr.Timeout() {
    t.Fatal("error is not a timeout")
  }
}
