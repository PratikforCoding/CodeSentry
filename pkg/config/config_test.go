package config_test

import (
	"os"
	"testing"

	"github.com/PratikforCoding/CodeSentry/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TestLoad_WithEnvVars tests that Load reads environment variables correctly.
func TestLoad_WithEnvVars(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("PORT", "9090")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("MONGO_URI", "mongodb://testuser:testpass@localhost:27017/testdb")

	defer func() {
		// Clean up environment variables after test
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("MONGO_URI")
	}()

	cfg := config.Load()

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, "mongodb://testuser:testpass@localhost:27017/testdb", cfg.MongoURI)
}

// TestLoad_Defaults tests that Load returns default values when env vars are not set.
func TestLoad_Defaults(t *testing.T) {
	// Ensure environment variables are unset
	os.Unsetenv("PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("MONGO_URI")

	cfg := config.Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "mongodb://root:example@mongo:27017/codesentry?authSource=admin", cfg.MongoURI)
}
