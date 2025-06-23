package proxy

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Config represents the proxy configuration from the config.json file
type Config struct {
	Proxy struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"proxy"`
	OriginServers []OriginServer `json:"origin_servers"`
}

type OriginServer struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Timeout string `json:"timeout"`
}

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

// GetTimeout returns the parsed timeout duration
func (c *OriginServer) GetTimeout() time.Duration {
	if c.Timeout == "" {
		return 30 * time.Second // default timeout
	}

	duration, err := time.ParseDuration(c.Timeout)
	if err != nil {
		log.Printf("Invalid timeout format, using default: %v", err)
		return 30 * time.Second
	}

	return duration
}
