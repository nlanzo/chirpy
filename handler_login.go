package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nlanzo/chirpy/internal/auth"
	"github.com/nlanzo/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if err := auth.CheckPasswordHash(params.Password, user.PasswordHash); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expiresIn := time.Hour // 1 hour
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to make JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to make refresh token", err)
		return
	}
	expiresAt := time.Now().Add(time.Hour * 24 * 60) // 60 days
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token: token,
		RefreshToken: refreshToken,
	})
	
}