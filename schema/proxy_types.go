package schema

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// OriginServer represents a server that the proxy will forward requests to, defined in the config.json file
type OriginServer struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Timeout string `json:"timeout"`
	Client  *http.Client
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

// ProxyServer represents the proxy server, it contains the config, clients, urls, names and reverse proxy
type ProxyServer struct {
	ProxyConfig  *ProxyConfig
	ReverseProxy http.Handler
}

func (p *ProxyServer) Start() {
	serverAddr := p.ProxyConfig.Host + p.ProxyConfig.Port
	log.Printf("Starting proxy server on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, p.ReverseProxy))
}

func (p *ProxyServer) SelectBackend(req *http.Request) (int, error) {
	serverName := req.URL.Query().Get("desired_server")
	if serverName == "" {
		return -1, fmt.Errorf("missing desired_server parameter")
	}

	for i, originServer := range p.ProxyConfig.OriginServers {
		if originServer.Name == serverName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("server '%s' not found", serverName)
}
