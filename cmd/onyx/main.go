package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// Standard Caddy modules
	_ "github.com/caddyserver/caddy/v2/modules/standard"

	// External plugins
	_ "github.com/caddy-dns/ovh"
	_ "github.com/corazawaf/coraza-caddy/v2"
)

func main() {
	caddycmd.Main()
}
