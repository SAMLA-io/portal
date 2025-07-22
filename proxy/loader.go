package proxy

import (
	"encoding/json"
	"log"
	"os"
	"portal/types"
)

type Config = types.Config

// Config represents the proxy configuration from the config.json file

// LoadConfig reads configuration from a JSON file and environment variables
// Environment variables take precedence over file config
func LoadConfig(filename string) (*Config, error) {
	var config Config

	// First, try to load from file if it exists
	fileConfig, err := loadConfigFromFile(filename)
	if err != nil {
		log.Printf("Warning: Could not load config from file %s: %v", filename, err)
		// Continue with environment-only config
	}

	// Load from environment variables
	config.LoadFromEnv()

	// Merge with file config (environment takes precedence)
	config.MergeWithFile(fileConfig)

	// Set defaults if still empty
	if config.Proxy.Port == "" {
		config.Proxy.Port = ":8080"
	}
	if config.Proxy.Host == "" {
		config.Proxy.Host = "localhost"
	}

	return &config, nil
}

// loadConfigFromFile reads configuration from a JSON file
func loadConfigFromFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
