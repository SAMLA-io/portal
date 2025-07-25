package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"portal/auth"
	schema "portal/types"
	"time"
)

type ProxyServer = schema.ProxyServer

func NewProxyServer(configPath string) (*ProxyServer, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	proxyConfig := &config.Proxy

	// create http clients for all origin servers from the config
	for i := range proxyConfig.OriginServers {
		proxyConfig.OriginServers[i].Client = &http.Client{
			Timeout: proxyConfig.OriginServers[i].GetTimeout(),
		}
	}

	proxyServer := &ProxyServer{
		ProxyConfig: proxyConfig,
	}

	// reverse proxy handler
	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("Request: %s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)
		start := time.Now()

		// find the origin server index from the query string
		originServerIndex, err := proxyServer.SelectBackend(req)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(rw, err)
			duration := time.Since(start)
			log.Printf("Response: %s %s -> STATUS: %d completed in %v", req.Method, req.URL.Path, http.StatusBadRequest, duration)
			return
		}

		originServerURL, err := url.Parse(proxyServer.ProxyConfig.OriginServers[originServerIndex].URL)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(rw, err)
			duration := time.Since(start)
			log.Printf("Response: %s %s -> STATUS: %d completed in %v", req.Method, req.URL.Path, http.StatusBadRequest, duration)
			return
		}

		// set req Host, URL and Request URI to forward a request to the origin server
		req.Host = originServerURL.Host
		req.URL.Host = originServerURL.Host
		req.URL.Scheme = originServerURL.Scheme
		req.RequestURI = ""

		// send a request to the origin server using configured client
		originServerResponse, err := proxyServer.ProxyConfig.OriginServers[originServerIndex].Client.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			duration := time.Since(start)
			log.Printf("Response: %s %s -> STATUS: %d completed in %v", req.Method, req.URL.Path, http.StatusInternalServerError, duration)
			return
		}
		defer originServerResponse.Body.Close()

		// copy the response headers and response from the origin server to the client
		for key, values := range originServerResponse.Header {
			for _, value := range values {
				rw.Header().Add(key, value)
			}
		}

		rw.WriteHeader(originServerResponse.StatusCode)
		_, err = io.Copy(rw, originServerResponse.Body)
		if err != nil {
			log.Printf("Error copying response body: %v", err)
			return
		}

		duration := time.Since(start)
		log.Printf("Response: %s %s -> STATUS: %d completed in %v", req.Method, req.URL.Path, originServerResponse.StatusCode, duration)
	})

	// Create a mux to handle both health check and proxy routes
	mux := http.NewServeMux()

	// Health check endpoint (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Proxy routes (with auth)
	mux.Handle("/", auth.VerifyingMiddleware(reverseProxy))

	proxyServer.ReverseProxy = mux
	return proxyServer, nil
}

func init() {

}
