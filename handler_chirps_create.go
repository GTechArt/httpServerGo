package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GTechArt/httpServerGo/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	type returnVal struct {
		Id         string `json:"id"`
		Created_at string `json:"created_at"`
		Updated_at string `json:"updated_at"`
		Body       string `json:"body"`
		UserId     string `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, valid := validatingChirp(params.Body)
	if !valid {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	uuidUser, err := uuid.Parse(params.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalide User_id: %s", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: uuid.NullUUID{UUID: uuidUser, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp : %s", err)
	}

	respondWithJSON(w, http.StatusCreated, returnVal{
		Id:         chirp.ID.String(),
		Created_at: chirp.CreatedAt.String(),
		Updated_at: chirp.UpdatedAt.String(),
		Body:       chirp.Body,
		UserId:     chirp.UserID.UUID.String(),
	})

}

func validatingChirp(msg string) (string, bool) {
	const maxChirpLenght = 140

	if len(msg) > maxChirpLenght {
		return "", false
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(msg, " ")
	for i, word := range words {
		for _, profaneWord := range profaneWords {
			if strings.Contains(strings.ToLower(word), strings.ToLower(profaneWord)) {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " "), true
}
