package middleware

import (
	"context"
	"net/http"
	"quickZ/utils"
	"strings"
)

// Define custom types for context keys
type ContextKey string

const (
	// Define the keys you will use in the context
	UserIDKey   ContextKey = "userID"
	UserTypeKey ContextKey = "userType"
)

// AuthMiddleware checks if the user has the correct role (e.g., "user")
func AuthMiddleware(requiredType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Check the user type
			if claims["type"] != requiredType {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			// Set user info into context for downstream handlers
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims["user_id"])
			ctx = context.WithValue(ctx, UserTypeKey, claims["type"])
			r = r.WithContext(ctx)

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}
