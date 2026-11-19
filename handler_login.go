package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nikunj-26/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	ctx := r.Context()

	var params parameters

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding request body", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(ctx, params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
