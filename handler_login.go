package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GTechArt/httpServerGo/internal/auth"
	"github.com/GTechArt/httpServerGo/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnVal struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode user login", err)
		return
	}
	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email or password is empty", nil)
		return
	}

	user, err := cfg.db.GetUser(req.Context(), params.Email)
	if err != nil {
		// add a fake cheking hash password for emulated an existing user
		auth.CheckPasswordHash("FakeP4ssW0r)", user.HashedPassword)

		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Check Hash Password to verify matching with account
	isValid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't verify hash password ", err)
		return
	}
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Counldn't generate JWT", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateToken(req.Context(), database.CreateTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(1440 * time.Hour),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
	}

	respondWithJSON(w, http.StatusOK, returnVal{
		User: User{
			Id:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        jwt,
			RefreshToken: refreshToken,
		},
	})
}
