package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/auth"
	"github.com/cooperstandard/NetZero/internal/util"
)

func (cfg *APIConfig) AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil || token != cfg.AdminKey {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (cfg *APIConfig) UserAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, 401, "unable to locate auth Header", err)
			return
		}
		id, err := auth.ValidateJWT(token, cfg.TokenSecret)
		if err != nil {
			util.RespondWithError(w, 401, "unable to decode jwt", err)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		next.ServeHTTP(w, r)
	}

}
