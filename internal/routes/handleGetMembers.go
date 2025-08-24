package routes

import (
	"net/http"

	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleGetMembers(w http.ResponseWriter, r *http.Request) {
	groupID, err := uuid.Parse(r.PathValue("groupID"))
	if err != nil {
		util.RespondWithError(w, 500, "invalid group id provided", err)
		return
	}

	members, err := cfg.DB.GetUsersByGroup(r.Context(), uuid.NullUUID{UUID: groupID, Valid: true})

	if err != nil {
		util.RespondWithError(w, 500, "couldn't get group members", err)
		return
	}

	var users []User

	for _, member := range members {
		users = append(users, User{
			ID:    member.ID,
			Email: member.Email,
			Name:  member.Name.String,
		})
	}

	util.RespondWithJSON(w, 200, users)

}
