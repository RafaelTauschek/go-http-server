package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/RafaelTauschek/http-server/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type returnVals struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find a token", err)
		return
	}

	token, err := cfg.db.GetUserFromRefreshToken(context.Background(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't authorize token", err)
		return
	}

	if token.ExpiresAt.Compare(time.Now()) == -1 {
		respondWithError(w, http.StatusUnauthorized, "Couldn't authorize token", errors.New("token is expired"))
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Couldn't authorize token", errors.New("token is revoked"))
		return
	}

	jwtToken, err := auth.MakeJWT(token.UserID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: jwtToken,
	})
}
