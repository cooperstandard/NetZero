package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	//TODO: fix this
	type parameters struct {
		GroupID string `json:"group_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	balances, err := cfg.DB.GetBalanceForDebtorByGroup(r.Context(), database.GetBalanceForDebtorByGroupParams{
		GroupID: uuid.NullUUID{Valid: true, UUID: uuid.MustParse(params.GroupID)},
		UserID:  uuid.NullUUID{Valid: true, UUID: r.Context().Value(UserID{}).(uuid.UUID)},
	})

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't locate balance records", err)
		return
	}

	util.RespondWithJSON(w, 200, balances)
}
