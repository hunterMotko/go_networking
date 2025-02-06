package echo

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

// this does not work on windows, wsl, or mac
func TestEchoServerUnixPacket(t *testing.T) {
  dir, err := os.MkdirTemp("", "echo_unixpacket")
  if err != nil {
    t.Fatal(err)
  }
  defer func() {
    if rErr := os.RemoveAll(dir); rErr != nil {
      t.Error(rErr)
    }
  }()
  ctx, cancel := context.WithCancel(context.Background())
  socket := filepath.Join(dir, fmt.Sprintf("s%d.sock", os.Getpid()))
  rAddr, err := streamingEchoServer(ctx, "unixpacket", socket)
  if err != nil {
    t.Fatal(err)
  }
  defer cancel()

  err = os.Chmod(socket, os.ModeSocket|0622)
  if err != nil {
    t.Fatal(err)
  }
  conn, err := net.Dial("unixpacket", rAddr.String())
  if err != nil {
    t.Fatal(err)
  }
  defer func() { _ = conn.Close() }()

  msg := []byte("ping")
  for i := 0; i < 3; i++ {
    _, err := conn.Write(msg)
    if err != nil {
      t.Fatal(err)
    }
  }
  buf := make([]byte, 1024)
  for i := 0; i < 3; i++ {
    n, err := conn.Read(buf)
    if err != nil {
      t.Fatal(err)
    }
    if !bytes.Equal(msg, buf[:n]) {
      t.Fatalf("expected reply %q; actual reply %q\n", msg, buf[:n])
    }
  }
  for i := 0; i < 3; i++ {
    _, err := conn.Write(msg)
    if err != nil {
      t.Fatal(err)
    }
  }
  buf = make([]byte, 2)
  for i := 0; i < 3; i++ {
    n, err := conn.Read(buf)
    if err != nil {
      t.Fatal(err)
    }
    if !bytes.Equal(msg[:n], buf[:n]) {
      t.Fatalf("expected reply %q; actual reply %q\n", msg[:2], buf[:n])
    }
  }
}

