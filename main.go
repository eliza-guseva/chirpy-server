package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eliza-guseva/chirpy-server/handlers"
	"github.com/eliza-guseva/chirpy-server/internal/db"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
		JWTSecret: os.Getenv("JWT_SECRET"),
		PolkaKey: os.Getenv("POLKA_KEY"),
	}

	fileServer := cfg.MiddlewareMetricsInc(http.FileServer(http.Dir("./static")))

	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("GET /api/healthz", handlers.Health)
	mux.HandleFunc("GET /admin/metrics", cfg.FSHits)
	mux.HandleFunc("POST /admin/reset", cfg.ResetUsers)

	mux.HandleFunc("POST /api/users", cfg.CreateUser)
	mux.HandleFunc("PUT /api/users", cfg.RequireAuth(cfg.UpdateUser))

	mux.HandleFunc("POST /api/login", cfg.Login)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshJWT)
	mux.HandleFunc("POST /api/revoke", cfg.RevokeRT)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.UpradeUserPolka)

	mux.HandleFunc("GET /api/chirps", cfg.GetChirps)
	mux.HandleFunc("POST /api/chirps", cfg.RequireAuth(cfg.CreateChirp))
	mux.HandleFunc("GET /api/chirps/{id}", cfg.GetChirp)
	mux.HandleFunc("DELETE /api/chirps/{id}", cfg.RequireAuth(cfg.DeleteChirp))



	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("Serving on http://%s", addr)
	server.ListenAndServe()

}
