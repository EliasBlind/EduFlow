package config

import (
	"flag"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Env     Env           `yaml:"env" env-default:"local" validate:"required,oneof=local dev prod"`
	Storage StorageConfig `yaml:"storage" validate:"required"`
	GRPC    GRPCConfig    `yaml:"grpc" validate:"required"`
}

type StorageConfig struct {
	Host         string `yaml:"host" validate:"required,hostname_rfc1123|ip"`
	Port         int    `yaml:"port" validate:"required,gte=1,lte=65535"`
	User         string `yaml:"user" validate:"required"`
	Password     string `env:"STORAGE_PASSWORD" env-required:"true"`
	DBName       string `yaml:"db_name" validate:"required"`
	Sslmode      string `yaml:"sslmode" validate:"oneof=disable enable verify-full"`
	MaxOpenConns int    `yaml:"max_open_conns" validate:"required,min=1"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host" validate:"required"`
	Port    int           `yaml:"port" validate:"required,gte=1,lte=65535"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func (e Env) IsLocal() bool { return e == EnvLocal }
func (e Env) IsDev() bool   { return e == EnvDev }
func (e Env) IsProd() bool  { return e == EnvProd }

func MustLoad() *Config {
	_ = godotenv.Load()

	var cfg Config
	path := fetchConfigPath()
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	cfg.validateConfigDate()
	return &cfg
}

// Получение пути к файлу конфигурации через флаг,
// через переменные окружения или через стандартный путь
// Приоритет: flag > env > default
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("JOURNAL_CONFIG_PATH")
	}

	if res == "" {
		res = "configs/journal_service/config.yaml"
	}

	if _, err := os.Stat(res); os.IsNotExist(err) {
		panic("config file is not exist: " + res)
	}
	return res
}

func (cfg *Config) validateConfigDate() {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic("error config validate: " + err.Error())
	}

	if cfg.Env.IsProd() && cfg.Storage.Sslmode == "disable" {
		panic("sslmode 'disable' is not allowed in production")
	}

	if cfg.Env.IsProd() && cfg.GRPC.Timeout <= 1*time.Second {
		panic("grpc timeout is too short (min 1s)")
	}
}
