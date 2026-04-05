package config

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnv_IsLocal(t *testing.T) {
	t.Parallel()
	assert.True(t, EnvLocal.IsLocal())
	assert.False(t, EnvDev.IsLocal())
	assert.False(t, EnvProd.IsLocal())
}

func TestEnv_IsDev(t *testing.T) {
	t.Parallel()
	assert.False(t, EnvLocal.IsDev())
	assert.True(t, EnvDev.IsDev())
	assert.False(t, EnvProd.IsDev())
}

func TestEnv_IsProd(t *testing.T) {
	t.Parallel()
	assert.False(t, EnvLocal.IsProd())
	assert.False(t, EnvDev.IsProd())
	assert.True(t, EnvProd.IsProd())
}

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

func TestMustLoad_PanicOnReadError(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	os.Chmod(tmpFile.Name(), 0000)
	defer os.Chmod(tmpFile.Name(), 0644)

	os.Setenv("JOURNAL_CONFIG_PATH", tmpFile.Name())
	defer os.Unsetenv("JOURNAL_CONFIG_PATH")
	os.Setenv("STORAGE_PASSWORD", "secret")
	defer os.Unsetenv("STORAGE_PASSWORD")

	assert.Panics(t, func() {
		MustLoad()
	})
}

func TestMustLoad_PanicOnValidation(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	content := `
env: "prod"
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

	os.Setenv("JOURNAL_CONFIG_PATH", tmpFile.Name())
	defer os.Unsetenv("JOURNAL_CONFIG_PATH")
	os.Setenv("STORAGE_PASSWORD", "secret")
	defer os.Unsetenv("STORAGE_PASSWORD")

	assert.Panics(t, func() {
		MustLoad()
	})
}

func TestFetchConfigPath_Priority_Flag(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	os.Args = []string{"test", "-config", tmpFile.Name()}

	result := fetchConfigPath()
	assert.Equal(t, tmpFile.Name(), result)
}

func TestFetchConfigPath_Priority_Env(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	os.Args = []string{"test"}
	os.Setenv("JOURNAL_CONFIG_PATH", tmpFile.Name())
	defer os.Unsetenv("JOURNAL_CONFIG_PATH")

	result := fetchConfigPath()
	assert.Equal(t, tmpFile.Name(), result)
}

func TestFetchConfigPath_Priority_Default(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"test"}
	os.Unsetenv("JOURNAL_CONFIG_PATH")

	defaultPath := "configs/journal_service/config.yaml"
	err := os.MkdirAll(filepath.Dir(defaultPath), 0755)
	assert.NoError(t, err)
	defer os.RemoveAll("configs")

	_, err = os.Create(defaultPath)
	assert.NoError(t, err)

	result := fetchConfigPath()
	assert.Equal(t, defaultPath, result)
}

func TestFetchConfigPath_InvalidExtension(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"test"}
	os.Setenv("JOURNAL_CONFIG_PATH", "config.txt")
	defer os.Unsetenv("JOURNAL_CONFIG_PATH")

	assert.Panics(t, func() {
		fetchConfigPath()
	})
}

func TestFetchConfigPath_FileNotExist(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"test"}
	os.Setenv("JOURNAL_CONFIG_PATH", "/nonexistent/path/config.yaml")
	defer os.Unsetenv("JOURNAL_CONFIG_PATH")

	assert.Panics(t, func() {
		fetchConfigPath()
	})
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
				Storage: StorageConfig{Port: 5432},
			},
		},
		{
			name: "Invalid sslmode value",
			cfg: Config{
				Env: EnvLocal,
				Storage: StorageConfig{
					Host: "localhost", Port: 5432, User: "u", DBName: "d",
					Sslmode: "invalid", MaxOpenConns: 1, Password: "p",
				},
				GRPC: GRPCConfig{Host: "h", Port: 80, Timeout: 5 * time.Second},
			},
		},
		{
			name: "Invalid Env value",
			cfg: Config{
				Env: "invalid_env",
				Storage: StorageConfig{
					Host: "localhost", Port: 5432, User: "u", DBName: "d",
					Sslmode: "disable", MaxOpenConns: 1, Password: "p",
				},
				GRPC: GRPCConfig{Host: "h", Port: 80, Timeout: 5 * time.Second},
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

func TestValidateConfigDate_ValidConfigs(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "Local with all valid fields",
			cfg: Config{
				Env: EnvLocal,
				Storage: StorageConfig{
					Host: "localhost", Port: 5432, User: "u", DBName: "d",
					Sslmode: "disable", MaxOpenConns: 10, Password: "p",
				},
				GRPC: GRPCConfig{Host: "h", Port: 8080, Timeout: 5 * time.Second},
			},
		},
		{
			name: "Prod with valid SSL and timeout",
			cfg: Config{
				Env: EnvProd,
				Storage: StorageConfig{
					Host: "db.prod.local", Port: 5432, User: "u", DBName: "d",
					Sslmode: "verify-full", MaxOpenConns: 20, Password: "p",
				},
				GRPC: GRPCConfig{Host: "0.0.0.0", Port: 443, Timeout: 10 * time.Second},
			},
		},
		{
			name: "Dev environment",
			cfg: Config{
				Env: EnvDev,
				Storage: StorageConfig{
					Host: "127.0.0.1", Port: 5432, User: "dev", DBName: "dev_db",
					Sslmode: "enable", MaxOpenConns: 5, Password: "dev_pass",
				},
				GRPC: GRPCConfig{Host: "localhost", Port: 50051, Timeout: 3 * time.Second},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				tt.cfg.validateConfigDate()
			})
		})
	}
}
