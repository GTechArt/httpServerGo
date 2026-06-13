package main

import (
	"encoding/json"
	"net/http"

	"github.com/GTechArt/httpServerGo/internal/auth"
	"github.com/GTechArt/httpServerGo/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleSetChirpyRed(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	type returnValue struct {
		User
	}

	apiKey, err := auth.GetAPIkey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get apiKey", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Apikey doesn't match", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		_, err := cfg.db.SetChirpyRedStatus(req.Context(), database.SetChirpyRedStatusParams{
			ID:          params.Data.UserId,
			IsChirpyRed: true,
		})
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't find the user", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}
