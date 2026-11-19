package main

import "net/http"

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching chirps", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
