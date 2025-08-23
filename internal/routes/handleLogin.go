package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cooperstandard/NetZero/internal/auth"
	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/util"
)

func (cfg *APIConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		ExpiresInSecs int    `json:"expiresInSeconds"` // TODO: for testing, configure this in the environment
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "unable to find user with given email", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		util.RespondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	expireIn := time.Duration(params.ExpiresInSecs) * time.Second
	if params.ExpiresInSecs < 1 || params.ExpiresInSecs > (3600) {
		expireIn = time.Hour
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.TokenSecret, expireIn)
	if err != nil {
		util.RespondWithError(w, 500, "unable to form jwt", err)
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()
	cfg.DB.CreateToken(r.Context(), database.CreateTokenParams{
		Token:     refreshToken,
		Email:     user.Email,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	res := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        jwt,
		RefreshToken: refreshToken,
	}
	util.RespondWithJSON(w, 200, res)
}
