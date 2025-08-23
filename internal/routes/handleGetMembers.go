package routes

import (
	"net/http"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleGetMembers(w http.ResponseWriter, r *http.Request) {
	groupID, ok := r.Context().Value("groupID").(uuid.UUID)
	if !ok {
		util.RespondWithError(w, 500, "invalid userID", nil)
		return
	}
	users, err := cfg.DB.GetUsersByGroup(r.Context(), database.GetUsersByGroupParams{
		GroupID: uuid.NullUUID{
			UUID:  groupID,
			Valid: true,
		},
	})

	if err != nil {
		return
	}

}
