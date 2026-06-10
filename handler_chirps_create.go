package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GTechArt/httpServerGo/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type returnVal struct {
		Chirp
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

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: params.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp : %s", err)
	}

	respondWithJSON(w, http.StatusCreated, returnVal{
		Chirp: Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		},
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
