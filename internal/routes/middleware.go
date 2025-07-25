package routes

import (
	"net/http"

	"github.com/cooperstandard/NetZero/internal/auth"
	"github.com/cooperstandard/NetZero/internal/util"
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

func (cfg *ApiConfig) UserAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, 401, "unable to locate auth Header", err)
			return
		}
		_, err = auth.ValidateJWT(token, cfg.TokenSecret)
		if err != nil {
			util.RespondWithError(w, 401, "unable to decode jwt", err)
			return
		}
		next.ServeHTTP(w, r)
	}
}
