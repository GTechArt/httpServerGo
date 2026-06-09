package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleChirpsValidate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLenght = 140

	if len(params.Body) > maxChirpLenght {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleaned := getCleanedBody(params.Body, profaneWords)

	respondWithJSON(w, http.StatusOK, returnVal{CleanedBody: cleaned})
}

func getCleanedBody(body string, profaneWords []string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		for _, profaneWord := range profaneWords {
			if strings.Contains(strings.ToLower(word), strings.ToLower(profaneWord)) {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}
