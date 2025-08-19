package main

import (
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"database/sql"
	"os"
	"github.com/joho/godotenv"
	"net/http"
	"github.com/eliza-guseva/chirpy-server/handlers"
	"github.com/eliza-guseva/chirpy-server/internal/db"
)


func main() {
	godotenv.Load()
	dbPool, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := db.New(dbPool)
	defer dbPool.Close()

	mux := http.NewServeMux()
	addr := "localhost:8080"
	cfg := &handlers.APIConfig{
		DBQueries: dbQueries,
	}

	fileServer := cfg.MiddlewareMetricsInc(http.FileServer(http.Dir("./static")))

	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("GET /api/healthz", handlers.Health)
	mux.HandleFunc("GET /admin/metrics", cfg.FSHits)
	mux.HandleFunc("POST /admin/reset", cfg.ResetUsers)
	mux.HandleFunc("POST /api/chirps", cfg.CreateChirp)
	mux.HandleFunc("POST /api/users", cfg.CreateUser)
	mux.HandleFunc("GET /api/chirps", cfg.GetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.GetChirp)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("Serving on http://%s", addr)
	server.ListenAndServe()

}
