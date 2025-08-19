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
			createClaims(userID, expiresIn),
		)
		signedToken, err := token.SignedString([]byte(tokenSecret))
		if err != nil {
			slog.Error("Error signing token", "error", err)
			os.Exit(1)
		}
		return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil},			
		)
		if err != nil {
			slog.Error("Error parsing token", "error", err)
			os.Exit(1)
		}
		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			return uuid.Parse(claims.Subject)
		}
		return uuid.Nil, nil
}

func createClaims(userID uuid.UUID, expiresIn time.Duration) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}}
