package main

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/chirpy/internal/auth"
	"example.com/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding user parameters")
		return
	}
	const maxExpirationSeconds = 1 * 60 * 60

	userRecord, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, userRecord.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(userRecord.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}
	tokenString := auth.MakeRefreshToken()
	// default expiration time is 60 days
	tokenParams := database.CreateRefreshTokenParams{
		Token:     tokenString,
		UserID:    userRecord.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}
	refreshTokenRecord, err := cfg.dbQueries.CreateRefreshToken(r.Context(), tokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token")
		return
	}
	returnUser := user{
		ID:           userRecord.ID,
		CreatedAt:    userRecord.CreatedAt,
		UpdatedAt:    userRecord.UpdatedAt,
		Email:        userRecord.Email,
		Token:        token,
		RefreshToken: refreshTokenRecord.Token,
	}
	respondWithJSON(w, http.StatusOK, returnUser)
}
