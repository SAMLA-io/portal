package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ProxyServer struct {
	config       *Config
	clients      []*http.Client
	urls         []*url.URL
	names        []string
	reverseProxy http.HandlerFunc
}

func NewProxyServer(configPath string) (*ProxyServer, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// create clients, urls and names slices for all origin servers from the config
	originServersClients := make([]*http.Client, len(config.OriginServers))
	originServersURLs := make([]*url.URL, len(config.OriginServers))
	originServersNames := make([]string, len(config.OriginServers))
	for i, originServer := range config.OriginServers {
		originServerURL, err := url.Parse(originServer.URL)
		if err != nil {
			log.Fatalf("Invalid origin server URL in config: %v", err)
		}

		originServersClients[i] = &http.Client{
			Timeout: originServer.GetTimeout(),
		}

		originServersURLs[i] = originServerURL
		originServersNames[i] = originServer.Name
	}

	proxyServer := &ProxyServer{
		config:  config,
		clients: originServersClients,
		urls:    originServersURLs,
		names:   originServersNames,
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// find the origin server index from the query string
		originServerIndex, err := proxyServer.selectBackend(req)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(rw, err)
			return
		}

		// set req Host, URL and Request URI to forward a request to the origin server
		req.Host = originServersURLs[originServerIndex].Host
		req.URL.Host = originServersURLs[originServerIndex].Host
		req.URL.Scheme = originServersURLs[originServerIndex].Scheme
		req.RequestURI = ""

		// send a request to the origin server using configured client
		originServerResponse, err := originServersClients[originServerIndex].Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}
		defer originServerResponse.Body.Close()

		// copy the response from the origin server to the client
		rw.WriteHeader(originServerResponse.StatusCode)
		_, err = io.Copy(rw, originServerResponse.Body)
		if err != nil {
			log.Printf("Error copying response: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})

	proxyServer.reverseProxy = reverseProxy
	return proxyServer, nil
}

func (p *ProxyServer) Start() {
	serverAddr := p.config.Proxy.Host + p.config.Proxy.Port
	log.Printf("Starting proxy server on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, p.reverseProxy))
}

func (p *ProxyServer) selectBackend(req *http.Request) (int, error) {
	serverName := req.URL.Query().Get("desired_server")
	if serverName == "" {
		return -1, fmt.Errorf("missing desired_server parameter")
	}

	for i, name := range p.names {
		if name == serverName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("server '%s' not found", serverName)
}

func init() {
	server, err := NewProxyServer("proxy/config.json")
	if err != nil {
		log.Fatal(err)
	}
	server.Start()
}
