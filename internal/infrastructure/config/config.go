package config

import (
	"os"
)

// Simple struct to hold the config
type Config struct {
	// Port is the port the HTTP server will listen on.
	Port string
	// PortFile is the file to read the ports info from
	PortFile string
}

// NewConfig creates a new Config instance.
func NewConfig() *Config {
	portFile := os.Getenv("PORT_FILE")
	if portFile == "" {
		portFile = "../../data/ports.json"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	return &Config{
		Port:     httpPort,
		PortFile: portFile,
	}
}
