package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/auth"
	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/util"
)

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
		Name: sql.NullString{String: params.Name, Valid: true},
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
