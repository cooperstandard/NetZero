package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleJoinGroup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupName string `json:"group_name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	group, err := cfg.DB.GetGroupByName(r.Context(), params.GroupName)
	if err != nil {
		util.RespondWithError(w, 500, "group not found", err)
		return
	}

	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		util.RespondWithError(w, 500, "invalid userID", nil)
		return
	}

	_, err = cfg.DB.JoinGroup(r.Context(), database.JoinGroupParams{
		UserID:  uuid.NullUUID{UUID: userID, Valid: true},
		GroupID: uuid.NullUUID{Valid: true, UUID: group.ID},
	})
	if err != nil {
		util.RespondWithError(w, 500, "unable to join group", err)
		return
	}

	w.WriteHeader(204)
}
