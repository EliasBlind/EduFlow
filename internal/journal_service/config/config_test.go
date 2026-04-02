package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMustLoad_Success(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

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
`
	tmpFile, err := os.CreateTemp("", "config*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	tmpFile.Close()

	// Устанавливаем окружение
	os.Setenv("JOURNAL_CONFIG_PATH", tmpFile.Name())
	os.Setenv("STORAGE_PASSWORD", "secret")
	defer os.Unsetenv("JOURNAL_CONFIG_PATH")
	defer os.Unsetenv("STORAGE_PASSWORD")

	var cfg *Config
	assert.NotPanics(t, func() {
		cfg = MustLoad()
	})

	assert.Equal(t, EnvLocal, cfg.Env)
	assert.Equal(t, "secret", cfg.Storage.Password)
	assert.Equal(t, 5*time.Second, cfg.GRPC.Timeout)
}

func TestValidateConfigDate_Table(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "Prod with disabled SSL",
			cfg: Config{
				Env: EnvProd,
				Storage: StorageConfig{
					Host: "127.0.0.1", Port: 5432, User: "u", DBName: "d",
					Sslmode: "disable", MaxOpenConns: 1, Password: "p",
				},
				GRPC: GRPCConfig{Host: "h", Port: 80, Timeout: 5 * time.Second},
			},
		},
		{
			name: "Prod with too short timeout",
			cfg: Config{
				Env: EnvProd,
				Storage: StorageConfig{
					Host: "127.0.0.1", Port: 5432, User: "u", DBName: "d",
					Sslmode: "verify-full", MaxOpenConns: 1, Password: "p",
				},
				GRPC: GRPCConfig{Host: "h", Port: 80, Timeout: 500 * time.Millisecond},
			},
		},
		{
			name: "Invalid Port (too high)",
			cfg: Config{
				Env:     EnvLocal,
				Storage: StorageConfig{Port: 99999},
			},
		},
		{
			name: "Missing required field (Host)",
			cfg: Config{
				Env:     EnvLocal,
				Storage: StorageConfig{Port: 5432}, // Host is empty
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

func TestFetchConfigPath_InvalidExtension(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	os.Setenv("JOURNAL_CONFIG_PATH", "config.txt")
	defer os.Unsetenv("JOURNAL_CONFIG_PATH")

	assert.Panics(t, func() {
		fetchConfigPath()
	})
}
