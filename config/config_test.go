package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefConfig(t *testing.T) {

	cfg := DefaultAppConfig
	if !cfg.System.Debug {
		t.Error("Expected Debug to be true")
	}
	// TestDatabase configuration
	if cfg.Database.Type != "postgres" {
		t.Errorf("Expected Database.Type 'postgres', got '%s'", cfg.Database.Type)
	}

	if cfg.Database.Name != "amts" {
		t.Errorf("Expected Database.Name 'amts', got '%s'", cfg.Database.Name)
	}

	if cfg.Database.MaxConn != 100 {
		t.Errorf("Expected Database.MaxConn 100, got %d", cfg.Database.MaxConn)
	}

	// Test logger configuration
	if cfg.Logger.Mode != "development" {
		t.Errorf("Expected Logger.Mode 'development', got '%s'", cfg.Logger.Mode)
	}

	if !cfg.Logger.FileEnable {
		t.Error("Expected Logger.FileEnable to be true")
	}

}

func TestLoadConfigFromFile(t *testing.T) {

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test-config.yml")

	// CreateTestconfig file
	configContent := `
system:
  debug: false


database:
  type: postgres
  host: localhost
  port: 5434
  name: amts
  user: postgres
  passwd: testpass
  max_conn: 50
  idle_conn: 5
  debug: true

logger:
  mode: production
  file_enable: false
  filename: /tmp/test.log
`

	err := os.WriteFile(configFile, []byte(configContent), 0600) //nolint:gosec // G306: test file permissions
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Loadconfiguration
	cfg := LoadConfig(configFile)
	if cfg.System.Debug {
		t.Error("Expected Debug to be false")
	}

	if cfg.Database.Type != "postgres" {
		t.Errorf("Expected Database.Type 'postgres', got '%s'", cfg.Database.Type)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected Database.Host 'localhost', got '%s'", cfg.Database.Host)
	}

	if cfg.Database.Port != 5434 {
		t.Errorf("Expected Database.Port 5434, got %d", cfg.Database.Port)
	}

	if cfg.Database.Name != "amts" {
		t.Errorf("Expected Database.Name 'amts', got '%s'", cfg.Database.Name)
	}

	if cfg.Database.MaxConn != 50 {
		t.Errorf("Expected Database.MaxConn 50, got %d", cfg.Database.MaxConn)
	}

}
