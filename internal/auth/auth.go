// Package auth provides authentication functions
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)




func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
		return "", err
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
			return "", err
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
			return uuid.Nil, err
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


func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", fmt.Errorf("no Authorization header found")
	}
	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid Authorization header")
	}
	return parts[1], nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {return "", err}
	return hex.EncodeToString(key), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	slog.Info("Getting API key", "headers", headers)
	header := headers.Get("Authorization")
	slog.Info("Authorization header", "header", header)
	parts := strings.Split(header, " ")
	if len(parts) !=2 || strings.ToLower(parts[0]) != "apikey" {
		return "", fmt.Errorf("invalid Authorization header")
	}
	return parts[1], nil
}
