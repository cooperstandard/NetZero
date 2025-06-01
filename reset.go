package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "403 Forbidden", nil)
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.RemoveAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to remove RemoveAllUsers", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
