package main

import (
	"encoding/json"
	"net/http"

	"github.com/GTechArt/httpServerGo/internal/auth"
	"github.com/GTechArt/httpServerGo/internal/database"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnVal struct {
		User
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't found bearer token", err)
		return
	}

	uuidUser, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Counldn't decode json body", err)
		return
	}

	hashPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.UpdateEmailPassword(req.Context(), database.UpdateEmailPasswordParams{
		ID:             uuidUser,
		Email:          params.Email,
		HashedPassword: hashPassword,
	})

	respondWithJSON(w, http.StatusOK, returnVal{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			ChirpyRed: user.IsChirpyRed,
		},
	})
}
