package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // postgres driver -> needs it
)

func TestCreateUser(t *testing.T) {
	godotenv.Load("../../.env")
	ctx := context.Background()
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	queries := New(tx)
	ctx = context.Background()

	testUser := CreateUserParams{
		Email:    "test@example.com",
		HashedPassword: "123456",
	}

	_, err = queries.CreateUser(ctx, testUser)


	if err != nil {
		t.Errorf("Error creating user: %v", err)
		return
	}
	slog.Info("User created successfully")
	t.Logf("User created successfully")
}
