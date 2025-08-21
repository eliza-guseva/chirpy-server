package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync/atomic"
	"github.com/google/uuid"
	"github.com/eliza-guseva/chirpy-server/internal/db"
	"github.com/eliza-guseva/chirpy-server/internal/auth"
	"log/slog"
	"context"
)

type APIConfig struct {
	fileserverHits atomic.Int32
	DBQueries *db.Queries
	JWTSecret string
	PolkaKey string
}


func (cfg *APIConfig) MiddlewareMetricsInc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		handler.ServeHTTP(w, r)
	})
}


func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK\n")
}

func (cfg *APIConfig) FSHits(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("templates/admin_metrics.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	data := struct{
		Hits int32
	}{
		Hits: cfg.fileserverHits.Load(),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templ.Execute(w, data)
}

func (cfg *APIConfig) RequireAuth(handler http.HandlerFunc) http.HandlerFunc {
      return func(w http.ResponseWriter, r *http.Request) {
          userID := cfg.Authenticate(w, r)
          if userID == uuid.Nil {
          	return
		  }

          ctx := context.WithValue(r.Context(), "userID", userID)
          handler(w, r.WithContext(ctx))
      }
  }


func (cfg *APIConfig) Authenticate(w http.ResponseWriter, r *http.Request) uuid.UUID {
	bearerToken, err := auth.GetBearerToken(r.Header)
	slog.Info("Bearer token", "token", bearerToken)
	if err != nil {
		slog.Error("Error getting bearer token", "error", err)
		respondWithError(w, 401, "Unauthorized")
		return uuid.Nil
	}
	userID, err := auth.ValidateJWT(bearerToken, cfg.JWTSecret)
	slog.Info("UserID", "userID", userID)
	if err != nil {
		slog.Error("Error validating JWT", "error", err)
		respondWithError(w, 401, "Unauthorized")
		return uuid.Nil
	}
	return userID
}


func respondWithError(w http.ResponseWriter, code int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}


