package sso

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = `C:\Users\broadcast\Desktop\BajoJaj\newBajoJajoo`
	}

	//err := godotenv.Load(envPath)
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error loading .env file from path: %s", envPath)
	}

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting canvas service")
	log.Debug("debug messages are enabled")

	// ctx := context.TODO()

	// postgreSQLClient, err := postgresql.NewClient(ctx, 5, cfg.PostgresConfig)
	// if err != nil {
	// 	log.Error("failed to init storage", "error", err)
	// 	os.Exit(1)
	// }
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
