package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrNoAuthHeader = errors.New("couldn't find Authorization header")
	ErrNoToken      = errors.New("couldn't find token in Authorization header")
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	mySigningKey := []byte(tokenSecret)

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := jwtToken.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", ErrNoAuthHeader
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		return "", ErrNoToken
	}
	return token, nil
}
