package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func init() {
	// Load configuration from JSON file
	config, err := LoadConfig("proxy/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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
			Transport: &http.Transport{
				Proxy: http.ProxyURL(originServerURL),
			},
		}

		originServersURLs[i] = originServerURL
		originServersNames[i] = originServer.Name
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if config.Logging.Enabled {
			fmt.Printf("[reverse proxy server] received request at: %s\n", time.Now())
		}

		// find the origin server index from the query string
		var originServerIndex int = -1
		for i, originServerName := range originServersNames {
			if originServerName == req.URL.Query().Get("origin_server") {
				originServerIndex = i
				break
			}
		}

		if originServerIndex == -1 {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(rw, "Server not found", http.StatusBadRequest)
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

	// Start server with configured port
	serverAddr := config.Proxy.Host + config.Proxy.Port
	log.Printf("Starting proxy server on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, reverseProxy))
}
