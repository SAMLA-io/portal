package auth

import (
	"net/http"
	"strings"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

// VerifyingMiddleware is a middleware that verifies the passed JWT token and extracts the user ID from it to check user permissions
// Verifies the JWT via Clerk's API and checks if the user has the necessary permissions to access the resource (has bought the product)
func VerifyingMiddleware(next http.Handler) http.Handler {
	return clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.Header.Get("Authorization")
		if jwtToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")

		productID := r.URL.Query().Get("product_id")

		permissions, err := VerifyUserPermissions(jwtToken, productID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !permissions {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}))
}
