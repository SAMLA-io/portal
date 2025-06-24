package schema

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
