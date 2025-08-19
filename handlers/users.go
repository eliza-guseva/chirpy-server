package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"os"
	"time"
)

type UserIn struct {
	Email string `json:"email"`
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
		slog.Error("Error decoding request: %s", err)
		return
	}
	if ! strings.Contains(reqUser.Email, "@") {
		slog.Error("Email must contain @")
		respondWithError(w, 400, "Email must contain @")
		return
	}

	user, err := cfg.DBQueries.CreateUser(r.Context(), reqUser.Email)
	if err != nil {
		slog.Error("Error creating user %s: %s", reqUser.Email, err)
		respondWithError(w, 500, "Could not create user")
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
		slog.Error("Error resetting users: %s", err)
		respondWithError(w, 500, "Could not reset users")
		return
	}
	respondWithJSON(w, 200, map[string]string{"message": "Users reset, Hits reset"})
}

