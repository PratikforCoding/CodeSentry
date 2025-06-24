package database_test

import (
	"testing"

	"github.com/PratikforCoding/CodeSentry/internal/database"
	"github.com/PratikforCoding/CodeSentry/pkg/config"
)

func TestInitMongoDB(t *testing.T) {
	cfg := &config.Config{
		MongoURI: "mongodb://localhost:27017", // Use test MongoDB URI
	}

	err := database.InitMongoDB(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if database.Client == nil {
		t.Fatal("expected MongoDB client to be initialized")
	}
}
