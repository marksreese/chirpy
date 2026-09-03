package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"example.com/chirpy/internal/auth"
	"example.com/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
		return
	}
	chirpRecord, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	returnChirp := chirp{
		ID:        chirpRecord.ID,
		CreatedAt: chirpRecord.CreatedAt,
		UpdatedAt: chirpRecord.UpdatedAt,
		Body:      chirpRecord.Body,
		UserID:    chirpRecord.UserID,
	}
	respondWithJSON(w, http.StatusOK, returnChirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps := []chirp{}
	chirpRecords, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirps")
		return
	}
	for _, currChirp := range chirpRecords {
		newChirp := chirp{
			ID:        currChirp.ID,
			CreatedAt: currChirp.CreatedAt,
			UpdatedAt: currChirp.UpdatedAt,
			Body:      currChirp.Body,
			UserID:    currChirp.UserID,
		}
		chirps = append(chirps, newChirp)
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	params := chirp{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding chirp")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user unauthorized")
		return
	}

	validBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirpRecord, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   validBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp")
		return
	}

	newChirp := chirp{
		ID:        chirpRecord.ID,
		CreatedAt: chirpRecord.CreatedAt,
		UpdatedAt: chirpRecord.UpdatedAt,
		Body:      chirpRecord.Body,
		UserID:    chirpRecord.UserID,
	}

	respondWithJSON(w, http.StatusCreated, newChirp)
}

func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", errors.New("chirp is too long")
	}

	filteredChirp := filterBadWords(chirp)

	return filteredChirp, nil
}

func filterBadWords(block string) string {
	bannedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(block, " ")
	var newWords []string

	for _, word := range words {
		lcWord := strings.ToLower(word)
		if _, banned := bannedWords[lcWord]; banned {
			newWords = append(newWords, "****")
		} else {
			newWords = append(newWords, word)
		}
	}
	newBlock := strings.Join(newWords, " ")
	return newBlock
}
