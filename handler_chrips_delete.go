package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/RafaelTauschek/http-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("chirpID")

	if param == "" {
		respondWithError(w, http.StatusBadRequest, "No parameter provided", nil)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't authenticate token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), uuid.MustParse(param))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "UserId can't delete chrips from other users", errors.New("not allowed"))
		return
	}

	err = cfg.db.DeleteChirp(context.Background(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudln't delete chirp", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
