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

func (cfg *APIConfig) HandleSettleUp(w http.ResponseWriter, r *http.Request) {
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

	tx, err := cfg.DBConn.Begin()
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

	debtorID, _ := r.Context().Value(UserID{}).(uuid.UUID)

	debts, err := qtx.GetUnpaidDebtsByCreditorAndDebtor(r.Context(), database.GetUnpaidDebtsByCreditorAndDebtorParams{
		Debtor:   debtorID,
		Creditor: uuid.MustParse(params.CreditorID),
		GroupID:  uuid.MustParse(params.GroupID),
	})
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve debt records", err)
		return
	}

	for _, debt := range debts {
		_, err = qtx.PayDebts(r.Context(), debt)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "unable to mark debts as paid", err)
			return
		}
	}

	zeroNumeric, _ := util.Numeric{Dollars: 0, Cents: 0}.ValidateAndFormNumericString()

	_, err = qtx.UpdateBalance(r.Context(), database.UpdateBalanceParams{
		Balance:    zeroNumeric,
		UserID:     debtorID,
		GroupID:    uuid.MustParse(params.GroupID),
		CreditorID: uuid.MustParse(params.CreditorID),
	})
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "unable to update balance record", err)
		return
	}

	w.WriteHeader(204)
	tx.Commit()
}
