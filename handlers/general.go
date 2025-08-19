package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync/atomic"
	"github.com/eliza-guseva/chirpy-server/internal/db"
)

type APIConfig struct {
	fileserverHits atomic.Int32
	DBQueries *db.Queries
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


