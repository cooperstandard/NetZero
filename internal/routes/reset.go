package routes

import (
	"fmt"
	"net/http"

	"github.com/cooperstandard/NetZero/internal/util"
)

// TODO: this should only be available in dev
func (cfg *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n", r.Header.Get("Authorization"))
	err := cfg.DB.RemoveAllUsers(r.Context())

	if err != nil {
		util.RespondWithError(w, 500, "failed to remove all users", err)
		return
	}
	w.WriteHeader(204)
}
