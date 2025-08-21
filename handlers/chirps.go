package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/eliza-guseva/chirpy-server/internal/db"
	"github.com/google/uuid"
)

type ChirpIn struct {
	Body string `json:"body"`
	UserID string `json:"user_id"`
}

type ChirpOut struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (cfg *APIConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	
	decoder := json.NewDecoder(r.Body)
	reqChirp := ChirpIn{}
	err := decoder.Decode(&reqChirp)
	if err != nil {
		slog.Error("Error decoding request", "error", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if len(reqChirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	_, fixed := checkForProfane(reqChirp.Body)
	UserID, _ := r.Context().Value("userID").(uuid.UUID)
	chirp, err := cfg.DBQueries.CreateChirp(r.Context(), 
		db.CreateChirpParams{
			Body: fixed,
			UserID: UserID,
		})
	if err != nil {
		slog.Error("Error creating chirp", "error", err)
		respondWithError(w, 500, "Could not create chirp")
		return
	}
	chirpOut := ChirpOut{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      fixed,
		UserID:    chirp.UserID.String(),
	}
	respondWithJSON(w, 201, chirpOut)

}


func (cfg *APIConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DBQueries.GetChirps(r.Context())
	if err != nil {
		slog.Error("Error getting chirps", "error", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	var chirpOut []ChirpOut	
	for _, chirp := range chirps {
		chirpOut = append(chirpOut, ChirpOut{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		})
	}
	respondWithJSON(w, 200, chirpOut)
}

func (cfg *APIConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	slog.Info(r.PathValue("id"))
	chID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("Invalid UUID", "error", err)
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}
	slog.Info("Getting chirp", "id", chID)
	chirp, err := cfg.DBQueries.GetChirp(r.Context(), chID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "Chirp not found")
			return
		}
		slog.Error("Error getting the chirp", "error", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	

	chirpOut := ChirpOut{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
	respondWithJSON(w, 200, chirpOut)
}

func (cfg *APIConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	chID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("Invalid UUID", "error", err)
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}
	authUserID := r.Context().Value("userID").(uuid.UUID)
	chirp, err := cfg.DBQueries.GetChirp(r.Context(), chID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "Chirp not found")
			return
		}
		slog.Error("Error getting chirp", "error", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if chirp.UserID != authUserID {
		slog.Error("User ID mismatch", "authUserID", authUserID, "chirpUserID", chirp.UserID)
		respondWithError(w, 403, "Unauthorized")
		return
	}
	err = cfg.DBQueries.DeleteChirp(r.Context(), chID)
	if err != nil {
		slog.Error("Error deleting chirp", "error", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	w.WriteHeader(204)
}

// Helpers 

func checkForProfane(chirp string) (hasProfate bool, fixed string) {
	hasProfane := false
	fixed = chirp
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, badWord := range badWords {
		for {
			lowerChirp := strings.ToLower(fixed)
			if ! strings.Contains(lowerChirp, badWord) {
				break
			}
			bwStart := strings.Index(lowerChirp, badWord)
			bwEnd := bwStart + len(badWord)
			fixed = fixed[:bwStart] + "****" + fixed[bwEnd:]
			hasProfane = true
		}
	}
	return hasProfane, fixed
}


