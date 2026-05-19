package config

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	envutil "github.com/EliasBlind/EduFlow/pkg/env"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         envutil.Env      `yaml:"env" env-default:"local" validate:"required,oneof=local dev prod"`
	PostgresSql PostgresqlConfig `yaml:"postgres" validate:"required"`
	GRPC        GRPCConfig       `yaml:"grpc" validate:"required"`
	Redis       RedisConfig      `yaml:"redis" validate:"required"`
	Mailtrap    MailtrapConfig   `yaml:"mailtrap" validate:"required"`
	SpiceDB     SpiceConfig      `yaml:"spice" validate:"required"`
	Token       TokenConfig      `yaml:"token" validate:"required"`
	Locale      LocaleConfig     `yaml:"locale" validate:"required"`
	Journal     JournalConfig    `yaml:"journal"`
}

type PostgresqlConfig struct {
	Host         string `yaml:"host" validate:"required,hostname_rfc1123|ip"`
	Port         int    `env:"POSTGRES_PORT" env-required:"true"`
	User         string `env:"POSTGRES_USER" env-required:"true"`
	Password     string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName       string `env:"POSTGRES_DB_NAME" env-required:"true"`
	Sslmode      string `yaml:"sslmode" validate:"oneof=disable enable verify-full"`
	MaxOpenConns int    `yaml:"max_open_conns" validate:"required,min=1"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host" validate:"required"`
	Port    int           `yaml:"port" validate:"required,gte=1,lte=65535"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type RedisConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `env:"REDIS_PORT" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	Db       int    `yaml:"db"`
}

type MailtrapConfig struct {
	Host        string `yaml:"host" validate:"required"`
	Port        int    `env:"MAILTRAP_PORT" env-required:"true"`
	Username    string `yaml:"username"`
	Password    string `env:"MAILTRAP_PASSWORD"`
	SenderEmail string `yaml:"sender_email" validate:"required,email"`
	UseTls      bool   `yaml:"use_tls"`
}

type SpiceConfig struct {
	Host     string        `yaml:"host" validate:"required"`
	Port     int           `env:"SPICEDB_PORT_F" env-required:"true"`
	Password string        `env:"SPICEDB_PASSWORD"`
	UseSSL   bool          `yaml:"use_ssl"`
	Timeout  time.Duration `yaml:"timeout"`
}

type TokenConfig struct {
	SecretKey       string        `env:"SECRET_KEY"`
	AccessTokenTTL  time.Duration `yaml:"access_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_ttl"`
	VerificationTTL time.Duration `yaml:"verification_ttl"`
}

type LocaleConfig struct {
	DefaultLang string `yaml:"default_lang"`
	Path        string `yaml:"path"`
}

type JournalConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func MustLoad() *Config {
	envPathStr, configPathStr := fetchPaths()

	envFile, err := envutil.LoadEnv(envPathStr)
	if err != nil {
		panic("failed to load env: " + err.Error())
	}

	godotenv.Load(envFile.Path())

	var cfg Config
	if err := cleanenv.ReadConfig(configPathStr, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	cfg.mustValidateConfigDate()

	return &cfg
}

// Получение пути к файлу конфигурации и env через флаг,
// через переменные окружения или через стандартный путь
// Приоритет: flag > env > default
func fetchPaths() (string, string) {
	var envPath string
	var configPath string

	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.StringVar(&envPath, "env", "", "path to .env file")
	fs.StringVar(&configPath, "config", "", "path to config file")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	if envPath == "" {
		envPath = os.Getenv("SSO_ENV_PATH")
	}

	if envPath == "" {
		envPath = "configs/sso_service/sso.env"
	}

	if configPath == "" {
		configPath = os.Getenv("SSO_CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "configs/sso_service/config.yaml"
	}

	if filepath.Ext(envPath) != ".env" {
		panic("invalid env extension: " + envPath)
	}
	if ext := filepath.Ext(configPath); ext != ".yaml" && ext != ".yml" {
		panic("invalid config extension: " + configPath)
	}

	return envPath, configPath
}

func (cfg *Config) mustValidateConfigDate() {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic("error config validate: " + err.Error())
	}

	if cfg.Env.IsProd() && cfg.PostgresSql.Sslmode == "disable" {
		panic("sslmode 'disable' is not allowed in production")
	}

	if cfg.Env.IsProd() && cfg.GRPC.Timeout <= 1*time.Second {
		panic("grpc timeout is too short (min 1s)")
	}

	if cfg.Token.SecretKey == "" {
		panic("SECRET_KEY is empty")
	}
}
