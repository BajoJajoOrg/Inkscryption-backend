package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	delivery "github.com/BajoJajoOrg/Inkscryption-backend/canvas/delivery/http"
	"github.com/BajoJajoOrg/Inkscryption-backend/canvas/repo"
	"github.com/BajoJajoOrg/Inkscryption-backend/canvas/usecase"
	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/middleware/logger"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/aws"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/db/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	// wd, err := os.Getwd()
	// fmt.Print("Current dir", wd)

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

	ctx := context.TODO()

	postgreSQLClient, err := postgresql.NewClient(ctx, 5, cfg.PostgresConfig)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	postgresRepo := repo.NewRepository(postgreSQLClient)

	awsClient, err := aws.NewClient(cfg.AWSConfig, log)
	if err != nil {
		log.Error("aws connection failed")
	}

	awsRepo := repo.NewAWSRepository(awsClient, cfg.AWSConfig)

	usecase := usecase.New(postgresRepo, &awsRepo, log)

	// _ = repository

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Csrft"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	delivery := delivery.New(cfg, router, usecase, log)

	if err = delivery.MapHandlers(); err != nil {
		log.Error("failed to map handlers")
	}

	errs := make(chan error, 2)

	go func() {
		errs <- delivery.ListenAndServe(cfg.HTTPServer)
	}()

	err = <-errs
	if err != nil {
		// log.Error("WEerr %s", err.Error())
		fmt.Printf("WEerr %s", err.Error())
	}
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
