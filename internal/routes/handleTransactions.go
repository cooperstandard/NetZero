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

type addDebtResult struct {
	err          error
	recordedDebt database.Debt
	failedDebt   database.CreateDebtParams
}

func (cfg *APIConfig) HandleCreateTransactions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Transactions []struct {
			Debtor string `json:"debtor"`
			Amount struct {
				Dollars int `json:"dollars"`
				Cents   int `json:"cents"`
			} `json:"amount"`
		} `json:"transactions"`
		TransactionID string `json:"transaction_id"`
		Creditor      string `json:"creditor"`
		GroupID       string `json:"group_id"`
		Title         string `json:"title"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Transactions) > 50 || len(params.Transactions) == 0 {
		util.RespondWithError(w, 400, "please batch records into groups of 50 or fewer to prevent service slow downs", fmt.Errorf("invalid number of records: %d", len(params.Transactions)))
		return
	}

	var transactionID uuid.UUID

	if params.TransactionID != "" {
		transactionID, err = uuid.Parse(params.TransactionID)
		if err != nil {
			util.RespondWithError(w, 500, "supplied transaction ID is invalid, please create a new transaction", err)
			return
		}
	} else {

		transaction, err := cfg.DB.CreateTransaction(r.Context(), database.CreateTransactionParams{
			Title:       params.Title,
			Description: sql.NullString{},
			AuthorID:    uuid.MustParse(params.Creditor),
			GroupID:     uuid.MustParse(params.GroupID),
		})
		if err != nil {
			util.RespondWithError(w, 500, "unable to create transaction record", err)
			return
		}
		transactionID = transaction.ID
	}
	resultChan := make(chan addDebtResult)
	var failed []database.CreateDebtParams
	var succeeded []database.Debt

	for _, v := range params.Transactions {
		go func(cfg APIConfig, ctx context.Context, debt database.CreateDebtParams, resultChan chan<- addDebtResult) {
			debtRecord, err := cfg.DB.CreateDebt(ctx, debt)
			if err != nil {
				resultChan <- addDebtResult{recordedDebt: debtRecord}
			} else {
				resultChan <- addDebtResult{err: err, failedDebt: debt}
			}
		}(*cfg, r.Context(), database.CreateDebtParams{
			Amount:        fmt.Sprintf("%d.%d", v.Amount.Dollars, v.Amount.Cents), // TODO: validate the Amount
			TransactionID: transactionID,
			Debtor:        uuid.MustParse(v.Debtor),
			Creditor:      uuid.MustParse(params.Creditor),
		}, resultChan)
	}

	for range len(params.Transactions) {
		result := <-resultChan
		if result.err != nil {
			succeeded = append(succeeded, result.recordedDebt)
		} else {
			failed = append(failed, result.failedDebt)
		}
	}

	util.RespondWithJSON(w, 200, struct {
		FailedTransactions     []database.CreateDebtParams `json:"failed_transactions,omitempty"`
		TransactionID          uuid.UUID                   `json:"transaction_id"`
		SuccessfulTransactions []database.Debt             `json:"successful_transactions,omitempty"`
	}{failed, transactionID, succeeded})

}

