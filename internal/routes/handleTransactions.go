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
	balance      database.Balance
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
				resultChan <- addDebtResult{err: err, failedDebt: debt}
			} else {
				balance, err := cfg.DB.InsertOrUpdateBalance(r.Context(), database.InsertOrUpdateBalanceParams{ // TODO: this and the earlier action should probably be part of the same transaction so we don't have to worry about intermediate errors
					Balance:    debt.Amount,
					UserID:     debt.Debtor,
					GroupID:    uuid.MustParse(params.GroupID), // TODO: pass this into the go func
					CreditorID: debt.Creditor,
				})
				if err != nil {
					cfg.DB.DeleteDebtById(r.Context(), debtRecord.ID)
					resultChan <- addDebtResult{err: err, failedDebt: debt}
					return
				}
				resultChan <- addDebtResult{recordedDebt: debtRecord, balance: balance}
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

func (cfg *APIConfig) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	// NOTE: expected to be either group_id= or author_id=, if neither is present assume getting transactions for the current user account
	groupID := r.URL.Query().Get("group_id")

	if groupID != "" {
		transactions, err := cfg.DB.GetTransactonsByAuthor(r.Context(), uuid.MustParse(groupID))
		if err != nil {
			util.RespondWithError(w, 404, "unable to locate records", err)
			return
		}
		util.RespondWithJSON(w, 200, transactions)
		return
	}

	authorID := r.URL.Query().Get("author_id")
	if authorID == "" {
		authorID = r.Context().Value(UserID{}).(uuid.UUID).String()
	}

	transactions, err := cfg.DB.GetTransactonsByAuthor(r.Context(), uuid.MustParse(authorID))
	if err != nil {
		util.RespondWithError(w, 404, "unable to locate records", err)
		return
	}
	util.RespondWithJSON(w, 200, transactions)
}

func (cfg *APIConfig) HandleGetTransactionDetails(w http.ResponseWriter, r *http.Request) {
	// NOTE: takes in a list from the query params and returns debts for those transaction ids
	// in the form of ?transactions=<transaction1id>,<transaction2id>,...,<transactionnid>
	transactions := r.URL.Query()["transactions"]
	if len(transactions) == 0 {
		w.WriteHeader(204)
		return
	}

	if len(transactions) > 100 {
		util.RespondWithError(w, 400, "please batch get transaction details requests into groups of size 100 or fewer", nil)
		return
	}

	transactionDetails := make(map[string][]database.Debt)

	for _, transaction := range transactions {
		transactionID := uuid.MustParse(transaction)
		debts, err := cfg.DB.GetDebtsByTransaction(r.Context(), transactionID)
		if err != nil {
			util.RespondWithError(w, 500, "unable to retrieve records for transaction: "+transaction, err)
			return
		}

		transactionDetails[transaction] = debts

	}

	util.RespondWithJSON(w, 200, transactionDetails)
}
