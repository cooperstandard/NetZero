package routes

import (
	"context"
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
			Debtor string `json:"debtor"`
			Amount struct {
				Dollars int `json:"dollars"`
				Cents   int `json:"cents"`
			} `json:"amount"`
		} `json:"transactions"`
		Creditor string `json:"creditor"`
		GroupID  string `json:"group_id"`
		Title    string `json:"title"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// TODO: this, need a way to add new debts to an existing transaction and limit the number of transactions which may be added in any one call to prevent locking up the database and web server
	// if len(params.Transactions) > 50 || len(params.Transactions) == 0 {
	// 	util.RespondWithError(w, 400, "please batch transactions into groups of 50 or fewer to prevent service slow downs", fmt.Errorf("invalid number of transactions: %d", len(params.Transactions)))
	// 	return
	// }

	transaction, err := cfg.DB.CreateTransaction(r.Context(), database.CreateTransactionParams{
		Title:       params.Title,
		Description: sql.NullString{},
		AuthorID:    uuid.MustParse(params.Creditor),
		GroupID:     uuid.MustParse(params.GroupID),
	})
	if err != nil {
		util.RespondWithError(w, 500, "unable to create transaction record", err)
	}

	okChan := make(chan bool)
	failedChan := make(chan database.CreateDebtParams)
	var failed []database.CreateDebtParams

	for _, v := range params.Transactions {
		go recordDebt(*cfg, r.Context(), database.CreateDebtParams{
			Amount:        fmt.Sprintf("%d.%d", v.Amount.Dollars, v.Amount.Cents), // TODO: validate the Amount
			TransactionID: transaction.ID,
			Debtor:        uuid.MustParse(v.Debtor),
			Creditor:      uuid.MustParse(params.Creditor),
		}, okChan, failedChan)
	}

	for range len(params.Transactions) {
		select {
		case <-okChan:
			continue
		case fail := <-failedChan:
			failed = append(failed, fail)
			continue
		}
	}

	if len(failed) > 0 {
		util.RespondWithJSON(w, 206, struct {
			FailedTransactions []database.CreateDebtParams `json:"failed_transactions"`
			TransactionID      uuid.UUID                   `json:"transaction_id"`
		}{failed, transaction.ID})
		return
	}

	util.RespondWithJSON(w, 200, struct {
		TransactionID uuid.UUID `json:"transaction_id"`
	}{transaction.ID})
}

func recordDebt(cfg APIConfig, ctx context.Context, debt database.CreateDebtParams, okChan chan bool, failedChan chan database.CreateDebtParams) {
	_, err := cfg.DB.CreateDebt(ctx, debt)
	if err != nil {
		failedChan <- debt
	}
}
