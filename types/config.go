package types

import (
	"os"
	"strings"
)

// Config represents the nucleus configuration from the config.json file
type Config struct {
	Proxy ProxyConfig `json:"proxy"`
}

// ProxyConfig represents the proxy configuration from the config.json file
type ProxyConfig struct {
	Port          string         `json:"port"`
	Host          string         `json:"host"`
	OriginServers []OriginServer `json:"origin_servers"`
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	// Load proxy config from environment
	if port := os.Getenv("PROXY_PORT"); port != "" {
		c.Proxy.Port = port
	}
	if host := os.Getenv("PROXY_HOST"); host != "" {
		c.Proxy.Host = host
	}

	// Load origin servers from environment
	// Format: ORIGIN_SERVERS=name1:url1:timeout1,name2:url2:timeout2
	if originServers := os.Getenv("ORIGIN_SERVERS"); originServers != "" {
		servers := strings.Split(originServers, ",")
		c.Proxy.OriginServers = make([]OriginServer, len(servers))

		for i, server := range servers {
			parts := strings.Split(server, ":")
			if len(parts) >= 2 {
				c.Proxy.OriginServers[i] = OriginServer{
					Name: parts[0],
					URL:  parts[1],
				}
				if len(parts) >= 3 {
					c.Proxy.OriginServers[i].Timeout = parts[2]
				} else {
					c.Proxy.OriginServers[i].Timeout = "30s" // default timeout
				}
			}
		}
	}
}

// MergeWithFile merges environment variables with file config
// Environment variables take precedence over file config
func (c *Config) MergeWithFile(fileConfig *Config) {
	if fileConfig == nil {
		return
	}

	// If environment doesn't override, use file config
	if c.Proxy.Port == "" {
		c.Proxy.Port = fileConfig.Proxy.Port
	}
	if c.Proxy.Host == "" {
		c.Proxy.Host = fileConfig.Proxy.Host
	}
	if len(c.Proxy.OriginServers) == 0 {
		c.Proxy.OriginServers = fileConfig.Proxy.OriginServers
	}
}
