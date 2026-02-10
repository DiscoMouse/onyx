// Package config handles the loading and persistence of Onyx configurations.
package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

// AdminConfig represents the root configuration for the onyx-admin tool.
type AdminConfig struct {
	Settings GlobalSettings `toml:"settings"`
	Nodes    []Node         `toml:"nodes"`
}

// GlobalSettings controls local application behavior.
type GlobalSettings struct {
	DefaultPort int    `toml:"default_port"`
	Theme       string `toml:"theme"`
}

// Node represents a paired Onyx engine.
type Node struct {
	Name     string    `toml:"name"`
	Address  string    `toml:"address"` // IP or Hostname
	Port     int       `toml:"port"`
	AddedAt  time.Time `toml:"added_at"`
	LastSeen time.Time `toml:"last_seen"`
}

// LoadConfig reads the TOML configuration from the specified path.
// If the file is missing, it returns a default configuration suitable for a fresh start.
func LoadConfig(path string) (*AdminConfig, error) {
	// Default state
	conf := &AdminConfig{
		Settings: GlobalSettings{
			DefaultPort: 2305,
			Theme:       "default",
		},
		Nodes: []Node{},
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return conf, nil
	}

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// SaveConfig writes the current configuration back to disk.
func (c *AdminConfig) SaveConfig(path string) error {
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

// AddNode safely appends or updates a node in the configuration.
func (c *AdminConfig) AddNode(name, address string, port int) {
	// Check for existing node to update
	for i, n := range c.Nodes {
		if n.Address == address && n.Port == port {
			c.Nodes[i].LastSeen = time.Now()
			c.Nodes[i].Name = name // Update name if changed
			return
		}
	}

	// Append new node
	c.Nodes = append(c.Nodes, Node{
		Name:     name,
		Address:  address,
		Port:     port,
		AddedAt:  time.Now(),
		LastSeen: time.Now(),
	})
}
