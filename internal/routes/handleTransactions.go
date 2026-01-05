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

	tx, err := cfg.DBConn.Begin()
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

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

		transaction, err := qtx.CreateTransaction(r.Context(), database.CreateTransactionParams{
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
	resultChan := make(chan addDebtResult, len(params.Transactions))
	var succeeded []database.Debt

	for _, v := range params.Transactions {
		// NOTE: was expirmenting with paralelizing this operation but need to look into how transactions work accross threads
		// go func(qtx *database.Queries, ctx context.Context, debt database.CreateDebtParams, resultChan chan<- addDebtResult) {
		func(qtx *database.Queries, ctx context.Context, debt database.CreateDebtParams, resultChan chan<- addDebtResult) {
			debtRecord, err := qtx.CreateDebt(ctx, debt)
			if err != nil {
				resultChan <- addDebtResult{err: err, failedDebt: debt}
			} else {
				balance, err := qtx.InsertOrUpdateBalance(r.Context(), database.InsertOrUpdateBalanceParams{
					Balance:    debt.Amount,
					UserID:     debt.Debtor,
					GroupID:    uuid.MustParse(params.GroupID), // TODO: pass this into the go func
					CreditorID: debt.Creditor,
				})
				if err != nil {

					resultChan <- addDebtResult{err: err, failedDebt: debt}
					return
				}
				resultChan <- addDebtResult{recordedDebt: debtRecord, balance: balance}
			}
		}(qtx, r.Context(), database.CreateDebtParams{
			Amount:        fmt.Sprintf("%d.%d", v.Amount.Dollars, v.Amount.Cents), // TODO: validate the Amount
			TransactionID: transactionID,
			Debtor:        uuid.MustParse(v.Debtor),
			Creditor:      uuid.MustParse(params.Creditor),
		}, resultChan)
	}

	for range len(params.Transactions) {
		result := <-resultChan
		if result.err == nil {
			succeeded = append(succeeded, result.recordedDebt)
		} else {
			util.RespondWithError(w, 405, "unable to create debt", nil)
			return
		}
	}

	tx.Commit()

	util.RespondWithJSON(w, 200, struct {
		TransactionID uuid.UUID       `json:"transaction_id"`
		Transactions  []database.Debt `json:"transactions,omitempty"`
	}{transactionID, succeeded})
}

func (cfg *APIConfig) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	// NOTE: expected to be either group_id= or author_id=, if neither is present assume getting transactions for the current user account
	groupID := r.URL.Query().Get("group_id")

	if groupID != "" {
		transactions, err := cfg.DB.GetTransactionsByGroup(r.Context(), uuid.MustParse(groupID))
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

func (cfg *APIConfig) HandleDeleteTransaction(w http.ResponseWriter, r *http.Request) {
	// TODO: do this in a database transaction
	type parameters struct {
		TransactionID string
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

	transaction, err := qtx.GetTransactionByID(r.Context(), uuid.MustParse(params.TransactionID))
	if err != nil {
		util.RespondWithError(w, 404, "unable to locate transaction record", err)
		return
	}

	debts, err := qtx.GetDebtsByTransaction(r.Context(), transaction.ID)
	if err != nil {
		util.RespondWithError(w, 404, "unable to locate individual debt records", err)
		return
	}

	for _, debt := range debts {
		if debt.Paid {
			_, err = qtx.InsertOrUpdateBalance(r.Context(), database.InsertOrUpdateBalanceParams{
				UserID:     debt.Creditor,
				GroupID:    transaction.GroupID,
				CreditorID: debt.Debtor,
				Balance:    debt.Amount,
			})
			if err != nil {
				util.RespondWithError(w, http.StatusInternalServerError, "internal server error", err)
				return
			}
		} else {
			balance, err := qtx.GetBalance(r.Context(), database.GetBalanceParams{
				UserID:     debt.Debtor,
				GroupID:    transaction.GroupID,
				CreditorID: debt.Creditor,
			})
			if err != nil {
				util.RespondWithError(w, http.StatusInternalServerError, "internal server error", err)
				return
			}
			newBalance := util.SimpleStringToNumeric(balance.Balance)
			newBalance, ok := newBalance.Subtraction(util.SimpleStringToNumeric(debt.Amount))
			balanceString, err := newBalance.ValidateAndFormNumericString()
			if !ok || err != nil {
				// TODO: this shouldn't happen because we have added correctly when storing the data
				w.WriteHeader(500)
				return
			}
			qtx.UpdateBalance(r.Context(), database.UpdateBalanceParams{
				Balance:    balanceString,
				UserID:     debt.Debtor,
				GroupID:    transaction.GroupID,
				CreditorID: debt.Creditor,
			})
			if err != nil {
				util.RespondWithError(w, http.StatusInternalServerError, "internal server error", err)
				return
			}
		}
		qtx.DeleteDebtById(r.Context(), debt.ID)
	}

	// TODO: query and rest of this implementation
	_, err = qtx.DeleteTransactionById(r.Context(), transaction.ID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	w.WriteHeader(204)
	tx.Commit()
}
