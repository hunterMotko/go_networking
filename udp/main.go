package main

import (
	"flag"
	"log"
	"network/udp/tftp"
	"os"
)

var (
  address = flag.String("a", "127.0.0.1:69", "listen address")
  payload = flag.String("p", "payload.svg", "file to server clients")
)

func main() {
  flag.Parse()

  p, err := os.ReadFile(*payload)
  if err != nil {
    log.Fatal(err)
  }
  s := tftp.Server{Payload: p}
  log.Fatal(s.ListenAndServe(*address))
}
