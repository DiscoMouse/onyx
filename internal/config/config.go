// Package config handles the parsing and validation of the Onyx configuration files.
// It supports TOML as the primary configuration format for the administration console.
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Server defines the connection and authentication details for a remote Onyx engine.
type Server struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	Cert    string `toml:"cert"`
	Key     string `toml:"key"`
}

// Config represents the root structure of the servers.toml file, containing a map of server aliases.
type Config struct {
	Servers map[string]Server `toml:"servers"`
}

// LoadConfig reads a TOML file from the specified path and decodes it into a Config struct.
// It returns an error if the file is missing or contains invalid TOML syntax.
func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, fmt.Errorf("failed to decode toml: %w", err)
	}

	return &conf, nil
}
