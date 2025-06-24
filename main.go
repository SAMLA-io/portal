package main

import (
	"fmt"
	"log"
	"net/http"
	"nucleus/proxy"
	"os"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
)

// loggingMiddleware logs each request with method, path, and response time
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the response time
		duration := time.Since(start)
		log.Printf("Response: %s %s completed in %v", r.Method, r.URL.Path, duration)
	})
}

func setup() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	clerk_secret_key := os.Getenv("CLERK_SECRET_KEY")
	clerk.SetKey(clerk_secret_key)
}

func main() {
	setup()
	// start the API server in a goroutine
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome the first server!")
		})

		handler := loggingMiddleware(mux)

		log.Println("API Server starting on :8081")
		log.Fatal(http.ListenAndServe(":8081", handler))
	}()

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome the second server!")
		})

		handler := loggingMiddleware(mux)

		log.Println("API Server starting on :8082")
		log.Fatal(http.ListenAndServe(":8082", handler))
	}()

	// start the proxy server in the main goroutine
	server, err := proxy.NewProxyServer("config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Proxy Server starting...")
	server.Start()
}
