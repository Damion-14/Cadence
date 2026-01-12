package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(r.Context())
				fmt.Printf("[ERROR] Request %s panicked: %v\n", requestID, err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]string{
						"code":       "INTERNAL_ERROR",
						"message":    "Internal server error",
						"request_id": requestID,
					},
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
