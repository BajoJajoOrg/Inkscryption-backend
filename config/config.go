package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer     `yaml:"http_server"`
	PostgresConfig `yaml:"postgres_config"`
	AWSConfig      `yaml:"aws_config"`
}

type HTTPServer struct {
	Host         string        `yaml:"host" env-default:"localhost"`
	Port         string        `yaml:"port" env-default:"6000"`
	Timeout      time.Duration `yaml:"timeout" env-default:"4s"`
	Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Database string `yaml:"database"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
}

type AWSConfig struct {
	AccessKey      string `yaml:"aws_access_key" env:"AWS_ACCESS_KEY_ID"`
	SecretKey      string `yaml:"aws_secret_key" env:"AWS_SECRET_ACCESS_KEY"`
	DefaultRegion  string `yaml:"aws_default_region" env:"AWS_DEFAULT_REGION"`
	BucketName     string `yaml:"bucket_name" env:"BUCKET_NAME"`
	SecretEndpoint string `env:"AWS_ENDPOINT"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg

}
