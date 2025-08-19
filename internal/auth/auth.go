// Package auth provides authentication functions
package auth

import (
	"log/slog"
	"os"
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)


func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
		os.Exit(1)
	}
	return string(hash), nil
}

func CheckPasswordHash(password string, hash string) error {
	slog.Info("Checking password", "password", password, "hash", hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil { return err }
	return nil
}

func MakeJWT(
	userID uuid.UUID, 
	tokenSecret string, 
	expiresIn time.Duration) (string, error) {
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.RegisteredClaims{
				Issuer: "chirpy",
				IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
				Subject: userID.String(),
			},
		)
		signedToken, err := token.SignedString([]byte(tokenSecret))
		if err != nil {
			slog.Error("Error signing token", "error", err)
			os.Exit(1)
		}
		return signedToken, nil
}
