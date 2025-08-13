package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/util"
)

func (cfg *ApiConfig) HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	group, err := cfg.DB.CreateGroup(r.Context(), params.Name)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "unable to create group", err)
	}

	ret := Group{
		Name:      group.Name,
		CreateAt:  group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
		ID:        group.ID,
	}

	util.RespondWithJSON(w, 200, ret)


}
