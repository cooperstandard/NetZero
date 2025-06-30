package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerOneChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, 500, "unable to parse input id param", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "", err)
		return
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Body      string    `json:"body"`
		UpdatedAt time.Time `json:"updated_at"`
		UserID    uuid.UUID `json:"user_id"`
	}

	ret := returnVals{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		Body:      chirp.Body,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, 200, ret)
}
