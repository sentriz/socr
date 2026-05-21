package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"

	"go.senan.xyz/socr/server/auth"
	"go.senan.xyz/socr/server/resp"
)

func WithCORS() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"DNT", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Range"}),
		handlers.MaxAge(1728000),
	)
}

func WithJWT(hmacSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if checkJWT(hmacSecret, r) || checkJWTParam(hmacSecret, r) {
				next.ServeHTTP(w, r)
				return
			}
			resp.Errorf(w, http.StatusUnauthorized, "unauthorised")
		})
	}
}

func WithJWTOrAPIKey(hmacSecret, apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if checkAPIKey(apiKey, r) || checkJWT(hmacSecret, r) || checkJWTParam(hmacSecret, r) {
				next.ServeHTTP(w, r)
				return
			}
			resp.Errorf(w, http.StatusUnauthorized, "unauthorised")
		})
	}
}

func WithLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("req %q", r.URL) //nolint:gosec
			next.ServeHTTP(w, r)
		})
	}
}

func checkAPIKey(apiKey string, r *http.Request) bool {
	header := r.Header.Get("X-Api-Key")
	return apiKey != "" && header == apiKey
}

func checkJWT(hmacSecret string, r *http.Request) bool {
	header := r.Header.Get("Authorization")
	header = strings.TrimPrefix(header, "bearer ")
	header = strings.TrimPrefix(header, "Bearer ")
	return auth.TokenParse(hmacSecret, header) == nil
}

func checkJWTParam(hmacSecret string, r *http.Request) bool {
	param := r.URL.Query().Get("token")
	return auth.TokenParse(hmacSecret, param) == nil
}
