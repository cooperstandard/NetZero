package routes

import (
	"net/http"

	"github.com/cooperstandard/NetZero/internal/util"
)

func (cfg *APIConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	err := cfg.DB.RemoveAllUsers(r.Context())

	if err != nil {
		util.RespondWithError(w, 500, "failed to remove all users", err)
		return
	}

	err = cfg.DB.RemoveAllGroups(r.Context())

	if err != nil {
		util.RespondWithError(w, 500, "failed to remove all groups", err)
		return
	}

	w.WriteHeader(204)
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
