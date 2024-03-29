package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"io"
	"log"
	"os"
	"time"
)

const (
	LocalEnv = "local"
	ProdEnv  = "prod"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := getConfigPath()
	if configPath == "" {
		panic("config path not provided")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("specified config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	return &cfg
}

func getConfigPath() string {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	path := flagSet.String("config", "", "")

	_ = flagSet.Parse(os.Args[1:])

	if path == nil || *path == "" {
		return os.Getenv("CONFIG_PATH")
	}

	return *path
}
