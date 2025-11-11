package config

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var authToken *jwtauth.JWTAuth

func GetToken(log *slog.Logger) *jwtauth.JWTAuth {
	if authToken == nil {
		key := os.Getenv("jwt_key")
		if key == "" {
			panic("jwt_key is not set")
		}
		authToken = jwtauth.New("HS256", key, nil)
	}
	return authToken
}

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	//StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"0s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"0s"`
	//User        string        `yaml:"user" env-required:"true"`
	//Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	var configPath string
	if runtime.GOOS == "windows" {
		configPath = "./config/local.yaml"
	} else {
		configPath = "./config/local.yaml"
	}
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
