package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GTechArt/httpServerGo/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresIn int    `json:"expires_in_seconds"`
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

	expriresIn := 1 * time.Hour
	if params.ExpiresIn != 0 {
		expriresIn = time.Duration(params.ExpiresIn) * time.Second
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expriresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Counldn't generate JWT Token", err)
	}

	respondWithJSON(w, http.StatusOK, returnVal{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
	})
}
