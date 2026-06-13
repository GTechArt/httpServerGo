package main

import (
	"net/http"

	"github.com/GTechArt/httpServerGo/internal/auth"
	"github.com/GTechArt/httpServerGo/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get UserID, bad token", err)
		return
	}

	chirpID := req.PathValue("chirpID")
	uuidChirp, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse chirp UUID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), uuidChirp)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if userId != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "You have not authorizate", nil)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), database.DeleteChirpParams{
		ID:     uuidChirp,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
