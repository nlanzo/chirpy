package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type cleanedChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if params.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Body is required", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}


	cleaned_chirp := replaceBadWords(params.Body)

	validatedChirp := cleanedChirp{
		CleanedBody: cleaned_chirp,
	}
	respondWithJSON(w, http.StatusOK, validatedChirp)
}

func replaceBadWords(chirp string) string {
	bad_words := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(chirp, " ")

	for i, word := range words {
		for _, bad_word := range bad_words {
			if strings.ToLower(word) == bad_word {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}