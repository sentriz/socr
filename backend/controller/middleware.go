package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"go.senan.xyz/socr/controller/auth"
)

func (c *Controller) WithCORS() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"DNT", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Range"}),
		handlers.MaxAge(1728000),
	)
}

func (c *Controller) WithAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			keyHeader := r.Header.Get("x-api-key")
			if c.APIKey != "" && keyHeader == c.APIKey {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("authorization")
			authHeader = strings.TrimPrefix(authHeader, "bearer ")
			authHeader = strings.TrimPrefix(authHeader, "Bearer ")
			if err := auth.TokenParse(c.HMACSecret, authHeader); err == nil {
				next.ServeHTTP(w, r)
				return
			}

			authParam := r.URL.Query().Get("token")
			if err := auth.TokenParse(c.HMACSecret, authParam); err == nil {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "unauthorised", http.StatusUnauthorized)
		})
	}
}

func (c *Controller) WithLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("req %q", r.URL)
			next.ServeHTTP(w, r)
		})
	}
}
