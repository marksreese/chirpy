package main

import (
	"net/http"

	"example.com/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user unauthorized")
		return
	}
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
