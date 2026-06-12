package main

import (
	"net/http"
	"time"

	"github.com/GTechArt/httpServerGo/internal/auth"
)

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, req *http.Request) {
	type returnVal struct {
		Token string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't found bearer token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Resquest Unothorized", err)
		return
	}

	if user.RevokedAt.Valid && time.Now().UTC().Compare(user.RevokedAt.Time) != -1 {
		respondWithError(w, http.StatusUnauthorized, "Token Expired", err)
		return
	}

	token, err := auth.MakeJWT(user.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Counldn't generate JWT Token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVal{
		Token: token,
	})
}

func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find bearer token", err)
		return
	}

	refreshToken, err := cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil || !refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh_token in database", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
