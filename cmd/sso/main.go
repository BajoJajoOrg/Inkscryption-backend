package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/middleware/logger"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/db/postgresql"
	"github.com/BajoJajoOrg/Inkscryption-backend/sso/delivery"
	"github.com/BajoJajoOrg/Inkscryption-backend/sso/repo"
	"github.com/BajoJajoOrg/Inkscryption-backend/sso/usecase"
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
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = `C:\Users\broadcast\Desktop\BajoJaj\newBajoJajoo\.env`
	}

	//err := godotenv.Load(envPath)
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error loading .env file from path: %s", envPath)
	}

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting auth service")
	log.Debug("debug messages are enabled")

	log.Info("application stopped")

	ctx := context.TODO()

	postgreSQLClient, err := postgresql.NewClient(ctx, 5, cfg.PostgresConfig)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	userRepo := repo.NewRepository(postgreSQLClient)

	userUC := usecase.New(userRepo, log)

	router := chi.NewRouter()

	// var tokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecretKey), nil)

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "hooli-pishem.ru"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Csrft"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	// router.Use(jwtauth.Verifier(tokenAuth))
	// router.Use(jwtauth.Authenticator(tokenAuth))

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	userDelivery := delivery.New(cfg, router, userUC, log, cfg.JWTSecretKey)

	if err = userDelivery.MapHandlers(); err != nil { // CANVAS MAPPING
		log.Error("failed to map user handlers")
	}

	errs := make(chan error, 2)

	go func() {
		errs <- userDelivery.ListenAndServe(cfg.HTTPServer)
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
