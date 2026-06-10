package main

import (
	"encoding/json"
	"net/http"

	"github.com/GTechArt/httpServerGo/internal/auth"
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

	isValid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't verify hash password ", err)
		return
	}
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVal{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
