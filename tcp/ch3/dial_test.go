package ch3_test

import (
	"io"
	"net"
	"testing"
)

// TODO: This example closes the connection is this what this is suppose to do here or
// our we actually supposed to make it to the buffer read?
func TestDial(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer func() { done <- struct{}{} }()
		for {
			conn, err := listener.Accept()
			if err != nil {
        t.Logf("ACCEPT ERR: %v\n", err)
				return
			}
			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()
				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("recieved: %q\n", buf[:n])
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
    t.Fatalf("DIAL ERROR: %v\n", err)
	}
	conn.Close()
	<-done
	listener.Close()
	<-done
}
