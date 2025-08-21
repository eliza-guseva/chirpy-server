package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/eliza-guseva/chirpy-server/internal/auth"
	"github.com/eliza-guseva/chirpy-server/internal/db"
	"github.com/google/uuid"
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
	Token     string    `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed bool    `json:"is_chirpy_red"`
}


type PolkaIn struct {
	Event string `json:"event"`
	Data struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

// HANDLERS

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
		slog.Error("Error hashing passwoRd", "error", err)
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
		IsChirpyRed: user.IsChirpyRed,
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
	reqUser, err := getRequestUser(w, r)
	if err != nil { return }

	user, err := cfg.getUser(w, r, reqUser)
	if err != nil { return }
	
	jwtToken, err := createTokenWithExp(user.ID, w)
	if err != nil {
		slog.Error("Error creating token", "error", err)
		respondWithError(w, 500, "Could not create token")
		return
	}
	refreshToken, err := cfg.setRefreshToken(w, r, user)
	if err != nil { return }
	

	respondWithJSON(w, 200, UserOut{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwtToken,
		RefreshToken: refreshToken,
		IsChirpyRed: user.IsChirpyRed,
	})
}


func (cfg *APIConfig) RefreshJWT(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		slog.Error("Error getting bearer token", "error", err)
		respondWithError(w, 401, "Could not get bearer token")
		return
	}
	dbToken, err := cfg.DBQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		slog.Error("Error getting refresh token", "error", err)
		respondWithError(w, 401, "Could not get refresh token")
		return
	}
	jwtToken, err := createTokenWithExp(dbToken.UserID, w)
	if err != nil { return }
	respondWithJSON(w, 200, map[string]string{"token": jwtToken})
}

func (cfg *APIConfig) RevokeRT(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		slog.Error("Error getting bearer token", "error", err)
		respondWithError(w, 401, "Could not get bearer token")
		return
	}
	cfg.DBQueries.ExpireRefreshToken(r.Context(), refreshToken)
	w.WriteHeader(204)

}

func (cfg *APIConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	authUserID := r.Context().Value("userID").(uuid.UUID)
	slog.Info("Auth user ID", "authUserID", authUserID)
	reqUser, err := getRequestUser(w, r)
	slog.Info("Request user", "reqUser", reqUser)
	if err != nil { return }

	dbUser, err := cfg.DBQueries.GetUserByID(r.Context(), authUserID)
	if err != nil { 
		slog.Error("Error getting user", "error", err)
		respondWithError(w, 500, "Could not get user")
		return
	}
	if dbUser.ID != authUserID {
		slog.Error("User ID mismatch", "authUserID", authUserID, "dbUser", dbUser.ID)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	hashedPassword, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
		respondWithError(w, 500, "Could not hash password")
		return
	}
	
	user, err := cfg.DBQueries.UpdateUser(r.Context(), db.UpdateUserParams{
		ID: dbUser.ID,
		Email: reqUser.Email,
		HashedPassword: hashedPassword,
	})
	
	if err != nil {
		slog.Error("Error updating user", "error", err)
		respondWithError(w, 500, "Could not update user")
		return
	}
	
	respondWithJSON(w, 200, UserOut{
		ID: user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *APIConfig) UpradeUserPolka(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		slog.Error("Error getting API key", "error", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if apiKey != cfg.PolkaKey {
		slog.Error("API key mismatch", "apiKey", apiKey, "cfg.PolkaKey", cfg.PolkaKey)
		respondWithError(w, 401, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	event := PolkaIn{}
	err = decoder.Decode(&event)
	slog.Info("Event", "event", event)
	if 	err != nil {
		respondWithError(w, 400, "Invalid request body")
		return
	}
	if event.Event != "user.upgraded" {
		slog.Info("Event NOT user.upgraded", "event", event)
		w.WriteHeader(204)
		return
	}
	slog.Info("Event IS user.upgraded", "event", event)
	_, err = cfg.DBQueries.UpgradeUser(r.Context(),uuid.MustParse(event.Data.UserID))
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "User not found")
			return
		}
		slog.Error("Error upgrading user", "error", err)
		respondWithError(w, 500, "Could not upgrade user")
		return
	}
	w.WriteHeader(204)
	return
}

// HELPERS

func createTokenWithExp(userID uuid.UUID, w http.ResponseWriter) (string, error) {
	jwtToken, err := auth.MakeJWT(userID, os.Getenv("JWT_SECRET"), time.Hour)
	if err != nil {
		slog.Error("Error creating token", "error", err)
		respondWithError(w, 500, "Could not create token")
		return "", err
	}
	return jwtToken, nil
}

func getRequestUser(w http.ResponseWriter, r *http.Request) (UserIn, error) {
	decoder := json.NewDecoder(r.Body)
	reqUser := UserIn{}
	err := decoder.Decode(&reqUser)
	if err != nil {
		slog.Error("Error decoding request: %s", "error", err)
		respondWithError(w, 400, "Could not decode request")
		return UserIn{}, err	
	 }
	if ! strings.Contains(reqUser.Email, "@") {
		slog.Error("Email must contain @")
		respondWithError(w, 400, "Email must contain @")
		return UserIn{}, err
	}
	return reqUser, nil
}


func (cfg *APIConfig) getUser(
	w http.ResponseWriter, 
	r *http.Request, 
	reqUser UserIn,
) (db.User, error) {
	user, err := cfg.DBQueries.GetUser(r.Context(), reqUser.Email)
	slog.Info("User", "email", reqUser.Email, "user", user)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 401, "Incorrect email or password")
			return db.User{}, err
		}
		slog.Error("Error getting user", "error", err)
		respondWithError(w, 500, "Could not get user")
		return db.User{}, err
	}
	if err := auth.CheckPasswordHash(reqUser.Password, user.HashedPassword); err != nil { 
		slog.Error("Error checking password", "error", err)
		respondWithError(w, 401, "Incorrect email or password")
		return db.User{}, err
	}
	return user, nil
}


func (cfg *APIConfig) setRefreshToken(w http.ResponseWriter, r *http.Request, user db.User) (string, error) {
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		slog.Error("Error creating refresh token", "error", err)
		respondWithError(w, 500, "Could not create refresh token")
		return "", err	
	}
	_, err = cfg.DBQueries.CreateRefreshToken(
		r.Context(),
		db.CreateRefreshTokenParams{
			Token: refreshToken,
			UserID: user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		},
	)
	if err != nil {
		slog.Error("Error expiring refresh token", "error", err)
		respondWithError(w, 500, "Could not expire refresh token")
		return "", err
	}
	return refreshToken, nil
}
