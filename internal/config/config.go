// Package config handles the loading and persistence of Onyx configurations.
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// AdminConfig represents the client-side configuration (servers.toml).
type AdminConfig struct {
	Servers []ServerEntry `toml:"servers"`
}

// ServerEntry defines a single remote Onyx engine.
type ServerEntry struct {
	Name string `toml:"name"`
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}

// LoadConfig reads the TOML configuration from the specified path.
func LoadConfig(path string) (*AdminConfig, error) {
	var conf AdminConfig
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

// SaveConfig writes the current configuration back to disk.
func (c *AdminConfig) SaveConfig(path string) error {
	// Ensure the directory exists before writing
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(c)
}

// AddServer appends a new server to the config if it doesn't already exist.
func (c *AdminConfig) AddServer(name, ip string, port int) {
	for _, s := range c.Servers {
		if s.IP == ip {
			return // Server already exists
		}
	}
	c.Servers = append(c.Servers, ServerEntry{
		Name: name,
		IP:   ip,
		Port: port,
	})
}
