//Package auth provides authentication functions
package auth

import (
	"log/slog"
	"os"
	"golang.org/x/crypto/bcrypt"
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
	//hashedPassword, err := HashPassword(password)
	//slog.Info("Hashed password", "hashedPassword", hashedPassword)
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
