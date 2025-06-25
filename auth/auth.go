package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
)

// VerifyingMiddleware is a middleware that verifies the passed JWT token and extracts the user ID from it to check user permissions
func VerifyingMiddleware(next http.Handler) http.Handler {
	return clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := extractUserIDFromAuthHeader(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Printf("Authenticated user: %s", userID)
		next.ServeHTTP(w, r)
	}))
}

// extractUserIDFromAuthHeader extracts the user ID from the Authorization header
func extractUserIDFromAuthHeader(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Verify the JWT token and extract claims
	claims, err := clerkjwt.Verify(context.Background(), &clerkjwt.VerifyParams{
		Token: token,
	})
	if err != nil {
		return "", fmt.Errorf("failed to verify token: %v", err)
	}

	// Extract user ID from the subject claim
	userID := claims.RegisteredClaims.Subject
	if userID == "" {
		return "", fmt.Errorf("no user ID found in token")
	}

	return userID, nil
}
