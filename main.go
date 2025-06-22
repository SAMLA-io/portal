package main

import (
	"fmt"
	"log"
	"net/http"
	"nucleus/auth"
	"time"
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

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the Nucleus API!")
	})

	mux.HandleFunc("/users/create", auth.CreateUserHandler)

	handler := loggingMiddleware(mux)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
