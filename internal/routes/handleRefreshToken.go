package routes

import (
	"net/http"
	"time"

	"github.com/cooperstandard/NetZero/internal/auth"
	"github.com/cooperstandard/NetZero/internal/util"
)

func (cfg *APIConfig) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh")
	if err != nil {
		util.RespondWithError(w, 400, "unknown error occured", err)
	}

	refreshToken, err := cfg.DB.GetToken(r.Context(), refreshCookie.Value)
	if err != nil {
		util.RespondWithError(w, 401, "unable to retrieve refresh token record", err)
		return
	}
	if time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		util.RespondWithError(w, 401, "unable to retrieve refresh token record", err)
		return
	}

	jwt, err := auth.MakeJWT(refreshToken.UserID, cfg.TokenSecret, time.Hour)
	if err != nil {
		util.RespondWithError(w, 500, "unable to form JWT", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}
	util.RespondWithJSON(w, 200, response{Token: jwt})
}
