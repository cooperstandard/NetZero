package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleGetBalanceDebtor(w http.ResponseWriter, r *http.Request) {
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
		GroupID: uuid.MustParse(params.GroupID),
		UserID:  r.Context().Value(UserID{}).(uuid.UUID),
	})

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't locate balance records", err)
		return
	}

	util.RespondWithJSON(w, 200, balances)
}

func (cfg *APIConfig) HandleGetBalanceCreditor(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupID    string `json:"group_id"`
		CreditorID string `json:"creditor_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	balances, err := cfg.DB.GetBalanceForCreditorByGroup(r.Context(), database.GetBalanceForCreditorByGroupParams{
		GroupID:    uuid.MustParse(params.GroupID),
		CreditorID: uuid.MustParse(params.CreditorID),
	})

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't locate balance records", err)
		return
	}

	util.RespondWithJSON(w, 200, balances)
}
