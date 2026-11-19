package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	chirpId, err := uuid.Parse(r.PathValue("chirpId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing chirpId", err)
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
