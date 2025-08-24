package routes

import (
	"fmt"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/util"
)

// TODO: this should only be available in dev
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
