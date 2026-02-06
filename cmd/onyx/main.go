package main

import (
	"fmt"
	"os"

	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// Standard Caddy modules
	_ "github.com/caddyserver/caddy/v2/modules/standard"

	// External plugins
	_ "github.com/caddy-dns/ovh"
	_ "github.com/corazawaf/coraza-caddy/v2"
)

// DEFAULT: "dev"
// This will be overwritten by the build system (GitHub Actions)
// using -ldflags "-X 'main.version=v1.0.0'"
var version = "dev"

func main() {
	// Intercept the "version" command
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("Onyx Version:  %s\n", version)
		fmt.Print("Caddy Version: ")
	}

	caddycmd.Main()
}
