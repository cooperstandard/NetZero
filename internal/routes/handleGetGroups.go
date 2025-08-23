package routes

import (
	"net/http"
	"time"

	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleGetGroups(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	groups, err := cfg.DB.GetGroupsByUser(r.Context(), uuid.NullUUID{
		UUID:  uuid.MustParse(userID),
		Valid: true,
	})

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "unable to locate group records", err)
		return
	}

	type ret struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}

	var resp []ret

	for _, group := range groups {
		resp = append(resp, ret{
			ID:        group.ID,
			Name:      group.Name,
			CreatedAt: group.CreatedAt,
		})
	}

	util.RespondWithJSON(w, 200, resp)
}
