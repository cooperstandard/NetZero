package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "unable to retrieve chirps", err)
		return
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Body      string    `json:"body"`
		UpdatedAt time.Time `json:"updated_at"`
		UserID    uuid.UUID `json:"user_id"`
	}

	ret := []returnVals{}
	for _, chirp := range chirps {
		ret = append(ret, returnVals{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			Body:      chirp.Body,
			UpdatedAt: chirp.UpdatedAt,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, 200, ret)
	return
}
