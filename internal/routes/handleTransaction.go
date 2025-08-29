package routes

import "net/http"

func (cfg *APIConfig) HandleCreateTransactions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Transactions []struct {
			Debtor   string `json:"debtor"`
			Creditor string `json:"creditor"`
			Amount   struct {
				Dollars int `json:"dollars"`
				Cents   int `json:"cents"`
			} `json:"amount"` //TODO: not sure if this is how amount should be modelled.
		} `json:"transactions"`
		GroupID string `json:"group_id"`
	}

}
