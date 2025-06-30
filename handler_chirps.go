package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cooperstandard/NetZero/internal/auth"
	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "unable to locate auth Header", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, 401, "unable to decode jwt", err)
		return
	}

	type parameters struct {
		Body   string    `json:"body"`
	}
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Body      string    `json:"body"`
		UpdatedAt time.Time `json:"updated_at"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't decode parameters", err)
		return
	}
	// if params.UserID != userID {
	// 	reqDump, _ := httputil.DumpRequest(r, true)
	//
	// 	fmt.Printf("REQUEST:\n%s", string(reqDump))
	// 	fmt.Println("----")
	// 	fmt.Printf("params: %s\n", params.UserID.String())
	// 	fmt.Printf("jwt:    %s\n", userID)
	// 	respondWithError(w, 401, "unauthenticated", nil)
	// 	return
	// }

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	ret := returnVals{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		Body:      chirp.Body,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, 201, ret)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
