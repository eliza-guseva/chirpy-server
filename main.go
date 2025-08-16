package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	addr := "localhost:8080"

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("Serving on http://%s", addr)
	server.ListenAndServe()

}
