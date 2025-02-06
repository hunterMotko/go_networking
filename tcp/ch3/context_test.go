package ch3_test

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

// Using context gives use more control over how we want to cancel the connection or why we want to cancel the connection
// using context is a built in way to handle this without doing all the checking ourselves when useing
// go concurrency or async functions
func TestDialContext(t * testing.T) {
  dl := time.Now().Add(5 * time.Second)
  ctx, cancel := context.WithDeadline(context.Background(), dl)
  defer cancel()

  var d net.Dialer
  d.Control = func(network, address string, c syscall.RawConn) error {
    // sleep long enough to reach the contexts deadline
    time.Sleep(5 * time.Second + time.Millisecond)
    return nil
  }

  conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
  if err != nil {
    conn.Close()
    t.Fatalf("DIAL CTX ERR: %v\n", err)
  }

  nErr, ok := err.(net.Error)
  if !ok {
    t.Errorf("not a net error?: %v\n", err)
  } else {
    if !nErr.Timeout() {
      t.Errorf("error is not a timeout: %v\n", err)
    }
  }
  if ctx.Err() != context.DeadlineExceeded {
    t.Errorf("expected dealine exceeded; actual: %v\n", ctx.Err())
  }
}
