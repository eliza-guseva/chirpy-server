package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/eliza-guseva/chirpy-server/handlers"

)

func TestHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/healthz", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Health)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
