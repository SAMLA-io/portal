package main

import (
	"log"
	"os"
	"portal/proxy"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	clerk_secret_key := os.Getenv("CLERK_SECRET_KEY")
	clerk.SetKey(clerk_secret_key)

	// Get config file path from environment or use default
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "config.json"
	}

	// start the proxy server in the main goroutine
	server, err := proxy.NewProxyServer(configPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Proxy Server starting...")
	server.Start()
}
