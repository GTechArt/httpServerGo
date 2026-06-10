package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {
	type returnVal struct {
		Chirp
	}

	dat, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps: %s", err)
		return
	}

	returnVals := make([]returnVal, len(dat))
	for i, chirp := range dat {
		returnVals[i] = returnVal{
			Chirp: Chirp{
				Id:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID,
			},
		}
	}

	respondWithJSON(w, http.StatusOK, returnVals)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {
	type returnVal struct {
		Chirp
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
		Chirp: Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		},
	})

}
