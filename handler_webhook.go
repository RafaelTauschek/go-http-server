package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RafaelTauschek/http-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Api key provided", err)
		return
	}

	if apiKey != cfg.apikey {
		respondWithError(w, http.StatusUnauthorized, "No Api key doesn't match", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameter", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, err = cfg.db.UpgradeUser(context.Background(), params.Data.UserID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
