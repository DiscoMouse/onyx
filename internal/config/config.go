package config

// Server defines the connection details for an Onyx instance
type Server struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	Cert    string `toml:"cert"`
	Key     string `toml:"key"`
}

// Config represents the full list of servers managed by the admin tool
type Config struct {
	Servers map[string]Server `toml:"servers"`
}
