package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cooperstandard/NetZero/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		ExpiresInSecs int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to find user with given email", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	expireIn := time.Duration(params.ExpiresInSecs) * time.Second
	if params.ExpiresInSecs < 1 || params.ExpiresInSecs > (3600) {
		expireIn = 3600 * time.Second
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.tokenSecret, expireIn)
	if err != nil {
		respondWithError(w, 500, "unable to form jwt", err)
		return
	}

	res := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwt,
	}
	respondWithJSON(w, 200, res)

}
