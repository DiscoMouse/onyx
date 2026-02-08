// Package main provides the entry point for the Onyx security engine.
// The engine acts as a specialized wrapper around the Caddy server,
// integrating custom WAF and DNS modules.
package main

import (
	"log"
	"os"

	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// Import standard Caddy modules and Onyx-specific plugins
	_ "github.com/caddy-dns/ovh"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/corazawaf/coraza-caddy/v2"
)

// main initializes and executes the Caddy command logic.
// This allows the Onyx engine to function as a background system service.
func main() {
	caddycmd.Main()
}

// init performs pre-flight checks before the engine starts.
// It handles environmental configuration and logging verbosity.
func init() {
	if os.Getenv("ONYX_VERBOSE") == "true" {
		log.Println("Onyx Security Engine: Initializing core modules...")
	}
}
