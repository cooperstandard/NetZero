package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleGetGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserID{}).(uuid.UUID)
	if !ok {
		util.RespondWithError(w, 500, "invalid userID", nil)
		return
	}

	groups, err := cfg.DB.GetGroupsByUser(r.Context(), uuid.NullUUID{
		UUID:  userID,
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

func (cfg *APIConfig) HandleGetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := cfg.DB.GetGroups(r.Context())
	if err != nil {
		util.RespondWithError(w, 500, "unable to get group", err)
		return
	}

	util.RespondWithJSON(w, 200, groups)
}

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

	userID, ok := r.Context().Value(UserID{}).(uuid.UUID)
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

func (cfg *APIConfig) HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	ret := Group{
		Name:      group.Name,
		CreateAt:  group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
		ID:        group.ID,
	}

	util.RespondWithJSON(w, 200, ret)

}
