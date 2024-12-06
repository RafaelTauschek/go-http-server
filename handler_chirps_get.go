package main

import (
	"context"
	"net/http"
	"sort"

	"github.com/RafaelTauschek/http-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getSortedChirps(ctx context.Context, id, sortDirection string) ([]database.Chirp, error) {
	var data []database.Chirp
	var err error

	if id != "" {
		data, err = cfg.db.GetChirpsByUser(ctx, uuid.MustParse(id))
	} else {
		data, err = cfg.db.GetChrips(ctx)
	}

	if err != nil {
		return nil, err
	}

	sortFunc := func(i, j int) bool {
		if sortDirection == "desc" {
			return data[i].CreatedAt.After(data[j].CreatedAt)
		}
		return data[i].CreatedAt.Before(data[j].CreatedAt)
	}

	sort.Slice(data, sortFunc)
	return data, nil
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	sortDirection := "asc"

	if sort == "desc" {
		sortDirection = "desc"
	}

	data, err := cfg.getSortedChirps(context.Background(), id, sortDirection)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chrips", err)
		return
	}

	var chrips []Chirp

	for _, chirp := range data {
		chrips = append(chrips, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chrips)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("chirpID")

	if param == "" {
		respondWithError(w, http.StatusBadRequest, "No parameter provided", nil)
		return
	}

	chirpID := uuid.MustParse(param)

	chirp, err := cfg.db.GetChirpById(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No chirp found", err)
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}
