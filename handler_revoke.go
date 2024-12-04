package main

import (
	"context"
	"net/http"

	"github.com/RafaelTauschek/http-server/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token provided", err)
		return
	}

	err = cfg.db.RevokeToken(context.Background(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke  token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
