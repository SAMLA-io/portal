package proxy

import (
	"encoding/json"
	"os"
	"portal/types"
)

type Config = types.Config

// Config represents the proxy configuration from the config.json file

// LoadConfig reads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
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
