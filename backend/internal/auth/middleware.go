package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/damion-14/cadence/backend/internal/middleware"
)

func Middleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, r, "Missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondUnauthorized(w, r, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]
			claims, err := ValidateToken(tokenString, secret)
			if err != nil {
				respondUnauthorized(w, r, "Invalid or expired token")
				return
			}

			ctx := middleware.SetUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func respondUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	requestID := middleware.GetRequestID(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":       "UNAUTHORIZED",
			"message":    message,
			"request_id": requestID,
		},
	})
}
