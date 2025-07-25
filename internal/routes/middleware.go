package routes

import (
	"net/http"

	"github.com/cooperstandard/NetZero/internal/auth"
)

func (cfg *ApiConfig) AdminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil || token != cfg.AdminKey {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (cfg *ApiConfig) userAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
	return next
}
