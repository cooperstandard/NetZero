package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleCreateTransactions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Transactions []struct {
			Debtor   string `json:"debtor"`
			Amount   struct {
				Dollars int `json:"dollars"`
				Cents   int `json:"cents"`
			} `json:"amount"`
		} `json:"transactions"`
		Creditor string `json:"creditor"`
		GroupID string `json:"group_id"`
		Title   string `json:"title"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Transactions) > 50 || len(params.Transactions) == 0 {
		util.RespondWithError(w, 400, "please batch transactions into groups of 50 or fewer to prevent service slow downs", fmt.Errorf("invalid number of transactions: %d", len(params.Transactions)))
		return
	}

	transaction, err := cfg.DB.CreateTransaction(r.Context(), database.CreateTransactionParams{
		Title:       params.Title,
		Description: sql.NullString{},
		AuthorID:    uuid.MustParse(params.Creditor),
		GroupID:     uuid.MustParse(params.GroupID),
	})
	
	if err != nil {
		util.RespondWithError(w, 500, "unable to create transaction record", err)
	}

	debts := []database.Debt{}
	for _, v := range params.Transactions { //TODO: this should use go routines and collect a slice of errors to send back with the successful transactions

		debt, _ := cfg.DB.CreateDebt(r.Context(), database.CreateDebtParams{
			Amount:        "",
			TransactionID: transaction.ID,
			Debtor:        uuid.MustParse(v.Debtor),
			Creditor:      uuid.MustParse(params.Creditor),
		})
		debts = append(debts, debt)
	}

	util.RespondWithJSON(w, 200, debts)

}
