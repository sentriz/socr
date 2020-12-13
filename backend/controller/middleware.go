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

func (c *Controller) WithJWT() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if checkJWT(c, r) || checkJWTParam(c, r) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "unauthorised", http.StatusUnauthorized)
		})
	}
}

func (c *Controller) WithAPIKey() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if checkAPIKey(c, r) {
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

func checkAPIKey(c *Controller, r *http.Request) bool {
	header := r.Header.Get("x-api-key")
	return c.APIKey != "" && header == c.APIKey
}

func checkJWT(c *Controller, r *http.Request) bool {
	header := r.Header.Get("authorization")
	header = strings.TrimPrefix(header, "bearer ")
	header = strings.TrimPrefix(header, "Bearer ")
	return auth.TokenParse(c.HMACSecret, header) == nil
}

func checkJWTParam(c *Controller, r *http.Request) bool {
	param := r.URL.Query().Get("token")
	return auth.TokenParse(c.HMACSecret, param) == nil
}
