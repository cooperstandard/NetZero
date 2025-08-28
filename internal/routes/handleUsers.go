package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func (cfg *APIConfig) HandleRegister(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "unable to hash password", err)
		return
	}
	userDetails := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
		Name:           sql.NullString{String: params.Name, Valid: true},
	}

	user, err := cfg.DB.CreateUser(r.Context(), userDetails)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Unable to create user record", err)
		return
	}
	ret := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Name:      user.Name.String,
	}
	util.RespondWithJSON(w, 201, ret)

}

func (cfg *APIConfig) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n", r.Header.Get("Authorization"))
	users, err := cfg.DB.GetUsers(r.Context())

	if err != nil {
		util.RespondWithError(w, 500, "failed to retrieve users", err)
		return
	}

	var ret []User
	for _, v := range users {
		ret = append(ret, User{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Email:     v.Email,
			Name:      v.Name.String,
		})
	}

	util.RespondWithJSON(w, 200, ret)

}

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
