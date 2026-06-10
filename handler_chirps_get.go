package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {
	type returnVal struct {
		Id         string `json:"id"`
		Created_at string `json:"created_at"`
		Updated_at string `json:"updated_at"`
		Body       string `json:"body"`
		UserId     string `json:"user_id"`
	}

	dat, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps: %s", err)
		return
	}

	returnVals := make([]returnVal, len(dat))
	for i, chirp := range dat {
		returnVals[i] = returnVal{
			Id:         chirp.ID.String(),
			Created_at: chirp.CreatedAt.String(),
			Updated_at: chirp.UpdatedAt.String(),
			Body:       chirp.Body,
			UserId:     chirp.UserID.UUID.String(),
		}
	}

	respondWithJSON(w, http.StatusOK, returnVals)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {
	type returnVal struct {
		Id         string `json:"id"`
		Created_at string `json:"created_at"`
		Updated_at string `json:"updated_at"`
		Body       string `json:"body"`
		UserId     string `json:"user_id"`
	}

	chirpID := req.PathValue("chirpID")
	uuidChirp, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse chirp UUID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), uuidChirp)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVal{
		Id:         chirp.ID.String(),
		Created_at: chirp.CreatedAt.String(),
		Updated_at: chirp.UpdatedAt.String(),
		Body:       chirp.Body,
		UserId:     chirp.UserID.UUID.String(),
	})

}
