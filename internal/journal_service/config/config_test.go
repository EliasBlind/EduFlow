package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustLoad_Success(t *testing.T) {
	content := `
env: "local"

storage:
  host: "localhost"
  port: 5432
  user: "user"
  db_name: "test_db"
  sslmode: "disable"
  max_open_conns: 10

grpc:
  host: "localhost"
  port: 50051
  timeout: 5s
`
	tmpFile, _ := os.CreateTemp("", "config*.yaml")
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(content)

	os.Setenv("JOURNAL_CONFIG_PATH", tmpFile.Name())
	os.Setenv("STORAGE_PASSWORD", "secret")
	defer os.Unsetenv("JOURNAL_CONFIG_PATH")
	defer os.Unsetenv("STORAGE_PASSWORD")

	assert.NotPanics(t, func() {
		cfg := MustLoad()
		assert.Equal(t, EnvLocal, cfg.Env)
		assert.Equal(t, 5432, cfg.Storage.Port)
		assert.Equal(t, "secret", cfg.Storage.Password)
	})
}

func TestValidateConfigDate_Panics(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
	}{
		{
			name: "Prod with disabled SSL",
			cfg: &Config{
				Env:     "prod",
				Storage: StorageConfig{Host: "localhost", Port: 5432, User: "u", DBName: "d", Sslmode: "disable", MaxOpenConns: 1},
			},
		},
		{
			name: "Invalid Port",
			cfg: &Config{
				Env:     "local",
				Storage: StorageConfig{Port: 99999}, // > 65535
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				tt.cfg.validateConfigDate()
			})
		})
	}
}
