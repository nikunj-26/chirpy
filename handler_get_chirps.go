package main

import (
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/nikunj-26/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error

	userIdString := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	if userIdString != "" {
		userId, parseErr := uuid.Parse(userIdString)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing userId", parseErr)
			return
		}

		dbChirps, err = cfg.db.GetChirpsForUser(r.Context(), userId)
	} else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	if sortOrder == "desc" {
		slices.SortFunc(chirps, func(a, b Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
