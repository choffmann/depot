package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestValidateItem(t *testing.T) {
	tests := []struct {
		name    string
		item    string
		stock   int
		wantErr bool
	}{
		{"valid", "spring", 40, false},
		{"empty stock", "spring", 0, false},
		{"upper bound", "spring", 10000, false},
		{"empty name", "", 40, true},
		{"whitespace name", "   ", 40, true},
		{"negative stock", "spring", -1, true},
		{"stock above limit", "spring", 10001, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateItem(tt.item, tt.stock)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateItem(%q, %d) = %v, want error: %v", tt.item, tt.stock, err, tt.wantErr)
			}
		})
	}
}

func TestMigrateIsRepeatable(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" {
		t.Fatal("DATABASE_URL is not set, start the database with just db")
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	for range 2 {
		if err := migrate(logger); err != nil {
			t.Fatalf("migrate: %v", err)
		}
	}

	pool, err := connect()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	var stock int
	err = pool.QueryRow(context.Background(), "SELECT stock FROM items WHERE name = 'washer'").Scan(&stock)
	if err != nil {
		t.Fatalf("seed item washer missing: %v", err)
	}
	if stock != 8 {
		t.Errorf("washer stock = %d, want 8", stock)
	}
}
