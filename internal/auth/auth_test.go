package auth

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "I'mCorrectP4ssword!"
	password2 := "meTooButDifferent<3"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = '%v', wantErr '%v'", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects '%v', got '%v'", tt.matchPassword, match)
			}
		})
	}
}

func TestJWTToken(t *testing.T) {
	userId := uuid.New()
	secretTest := "CHUT!-Is-my-little-secret."

	tests := []struct {
		name         string
		verifySecret string
		expiresIn    time.Duration
		wantErr      bool
	}{
		{
			name:         "Same secret",
			verifySecret: secretTest,
			expiresIn:    time.Hour,
			wantErr:      false,
		},
		{
			name:         "Secret doesn't match",
			verifySecret: "Another-secret",
			expiresIn:    time.Hour,
			wantErr:      true,
		},
		{
			name:         "token expired",
			verifySecret: secretTest,
			expiresIn:    -time.Hour,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tokenString, err := MakeJWT(userId, secretTest, tt.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT() internal error: %s", err)
			}

			claims := jwt.RegisteredClaims{}
			_, err = jwt.ParseWithClaims(tokenString, &claims, func(*jwt.Token) (interface{}, error) {
				return []byte(tt.verifySecret), nil
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = '%v', wantErr '%v'", err, tt.wantErr)
			}
			if !tt.wantErr && claims.Subject != userId.String() {
				t.Errorf("MakeJWT() expects = '%v', got '%v'", userId, claims)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name              string
		headerSetup       func() http.Header
		verifyTokenString string
		wantErr           error
	}{
		{
			name: "Good Token",
			headerSetup: func() http.Header {
				h := http.Header{}
				h.Set("Authorization", "Bearer it-is-my-little-token-secret")
				return h
			},
			verifyTokenString: "it-is-my-little-token-secret",
			wantErr:           nil,
		},
		{
			name: "No Auth Found",
			headerSetup: func() http.Header {
				h := http.Header{}
				return h
			},
			verifyTokenString: "",
			wantErr:           ErrNoAuthHeader,
		},
		{
			name: "No Token Found",
			headerSetup: func() http.Header {
				h := http.Header{}
				h.Set("Authorization", "Bearer ")
				return h
			},
			verifyTokenString: "",
			wantErr:           ErrNoToken,
		},
		{
			name: "Tricky Token (starts/ends with cutset chars)",
			headerSetup: func() http.Header {
				h := http.Header{}
				h.Set("Authorization", "Bearer rabbit-ate-a-bear")
				return h
			},
			verifyTokenString: "rabbit-ate-a-bear",
			wantErr:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tokenTest, err := GetBearerToken(tt.headerSetup())
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetBearerToken() error = '%v', want '%v'", err, tt.wantErr)
			}
			if tt.wantErr == nil && tokenTest != tt.verifyTokenString {
				t.Errorf("GetBearerToken() expects = '%v', got '%v'", tt.verifyTokenString, tokenTest)
			}
		})
	}
}
