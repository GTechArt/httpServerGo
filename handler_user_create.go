package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleAddUesrs(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type returnVal struct {
		Id         string `json:"id"`
		Created_at string `json:"created_at"`
		Updated_at string `json:"updated_at"`
		Email      string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode user", err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user : %s", err)
	}

	respondWithJSON(w, http.StatusCreated, returnVal{
		Id:         user.ID.String(),
		Created_at: user.CreatedAt.String(),
		Updated_at: user.UpdatedAt.String(),
		Email:      user.Email,
	})

}
