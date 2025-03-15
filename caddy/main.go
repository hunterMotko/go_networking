package main

import (
	cmd "github.com/caddyserver/caddy/v2/cmd"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	// Injecting custom modules into caddy
  _ "network/caddy/restrict_prefix"
  _ "network/caddy/toml_adapter"
)

func main() {
	cmd.Main()
}
