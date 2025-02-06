package ch4

import (
	"errors"
	"log"
	"net"
	"time"
)

var (
	err error
	n   int
	i   = 7 // max number of retries
)

func reconnection() error {
	conn, err := net.Dial("tcp", "127.0.0.1:")
	if err != nil {
		log.Fatal(err)
	}
	for ; i > 0; i-- {
		n, err = conn.Write([]byte("hello world"))
		if err != nil {
			if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
				log.Println("temp err:", nErr)
				time.Sleep(10 * time.Second)
        continue
			}
      return err
		}
    break
	}
  if i == 0 {
    return errors.New("temporary write failure threshold exceeded")
  }
  log.Printf("wrote %d bytes to %s\n", n, conn.RemoteAddr())
  return nil
}
