package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nlanzo/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "API key not found", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Invalid event", nil)
		return
	}

	user, err := cfg.db.GetUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	_, err = cfg.db.UserAddChirpyRed(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}