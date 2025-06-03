package middleware

import (
	"context"
	"net/http"
)

func UserAgentMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "user-agent", r.Header.Get("User-Agent"))
		next(w, r.WithContext(ctx))
	}
}
