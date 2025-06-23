package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwks keyfunc.Keyfunc

func VerifyingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		var tokenString string
		fmt.Sscanf(authHeader, "Bearer %s", &tokenString)

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Example: access user's Clerk ID
		sub := claims["sub"]
		fmt.Fprintf(w, "Hello user %v\n", sub)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	jwksURL := os.Getenv("CLERK_JWKS_URL")

	// Create the JWKS from the resource at the given URL.
	var jwksErr error
	jwks, jwksErr = keyfunc.NewDefault([]string{jwksURL})
	if jwksErr != nil {
		log.Fatalf("Failed to get JWKS: %v", jwksErr)
	}
}
