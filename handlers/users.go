package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"os"
	"time"
	"database/sql"
	"github.com/eliza-guseva/chirpy-server/internal/db"
	"github.com/eliza-guseva/chirpy-server/internal/auth"
)

type UserIn struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserOut struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
}

func (cfg *APIConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqUser := UserIn{}
	err := decoder.Decode(&reqUser)
	if err != nil {
		slog.Error("Error decoding request", "error", err)
		return
	}
	if ! strings.Contains(reqUser.Email, "@") {
		slog.Error("Email must contain @")
		respondWithError(w, 400, "Email must contain @")
		return
	}

	hashedPassword, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
		respondWithError(w, 500, "Could not hash password")
		return
	}
	
	user, err := cfg.DBQueries.CreateUser(
		r.Context(), 
		db.CreateUserParams{
			Email: reqUser.Email, 
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		slog.Error("Error creating user", "error", err, "email", reqUser.Email)
		respondWithError(w, 500, "Could not create user")
		return 
	}
	userOut := UserOut{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, 201, userOut)
}


func (cfg *APIConfig) ResetUsers(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.DBQueries.ResetUsers(r.Context())
	if err != nil {
		slog.Error("Error resetting users", "error", err)
		respondWithError(w, 500, "Could not reset users")
		return
	}
	respondWithJSON(w, 200, map[string]string{"message": "Users reset, Hits reset"})
}

func (cfg *APIConfig) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqUser := UserIn{}
	err := decoder.Decode(&reqUser)
	if err != nil {
		slog.Error("Error decoding request: %s", "error", err)
		respondWithError(w, 400, "Could not decode request")
		return
	}
	if ! strings.Contains(reqUser.Email, "@") {
		slog.Error("Email must contain @")
		respondWithError(w, 400, "Email must contain @")
		return
	}

	user, err := cfg.DBQueries.GetUser(r.Context(), reqUser.Email)
	slog.Info("User", "email", reqUser.Email, "user", user)
	slog.Info("Password", "password", reqUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 401, "Incorrect email or password")
			return
		}
		slog.Error("Error getting user", "error", err)
		respondWithError(w, 500, "Could not get user")
		return
	}
	if err := auth.CheckPasswordHash(reqUser.Password, user.HashedPassword); err != nil { 
		slog.Error("Error checking password", "error", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	respondWithJSON(w, 200, UserOut{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

