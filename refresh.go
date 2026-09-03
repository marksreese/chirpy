package main

import (
	"net/http"
	"time"

	"example.com/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user unauthorized")
		return
	}
	user, err := cfg.dbQueries.GetUserByRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user from token")
		return
	}
	tokenString, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't make JWT")
		return
	}
	returnToken := struct {
		Token string `json:"token"`
	}{Token: tokenString}
	respondWithJSON(w, http.StatusOK, returnToken)
}
