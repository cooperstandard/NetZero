package routes

import (
	"fmt"
	"net/http"
)

// TODO: this should only be available in dev
func (cfg *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n", r.Header.Get("Authorization"))
	err := cfg.DB.RemoveAllUsers(r.Context())
	w.WriteHeader(204)
}

