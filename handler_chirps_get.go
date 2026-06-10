package main

import (
	"net/http"
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
